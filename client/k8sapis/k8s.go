package k8sapis

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	// kube-config (outside of cluster)
	"flag"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	// kube-config (in the pod)
	"k8s.io/client-go/rest"

	//"log"
	"encoding/json"
	"strconv"

	"io/ioutil"

	//ffmt "gopkg.in/ffmt.v1"
	//"github.com/TylerBrock/colorjson"
	//"github.com/imdario/mergo"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema" // GetResourceDynamically
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"

	// GetResourcesByJq
	//"github.com/itchyny/gojq"
	//"k8s.io/apimachinery/pkg/runtime"

	rbacv1 "k8s.io/api/rbac/v1"
)

// global variables
//var kubeconfig *string
var kubeclient *kubernetes.Clientset
var dynamicinterface dynamic.Interface
var metricsclient *metrics.Clientset
var contextbg context.Context

func getClusterConfig() (*rest.Config, error) {
	//_1st_In_the_pod_______________________________________________________________________________________________
	// create the in-cluster config
	kubeconfig, err := rest.InClusterConfig()
	if err != nil {
		//_2nd_Outside_of_cluster_______________________________________________________________________________________
		configpath := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
		if configpath == "" {
			if home := homedir.HomeDir(); home != "" {
				configpath = *flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
			} else {
				configpath = *flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
			}
			flag.Parse()
		}

		// use the current context in kubeconfig
		kubeconfig, err = clientcmd.BuildConfigFromFlags("", configpath)
		if err != nil {
			return nil, err
		}
	}

	return kubeconfig, err
}

//func getClientset() (*kubernetes.Clientset, error) {
func getClientset() {
	kubeconfig, err := getClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	kubeclient, err = kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	dynamicinterface, err = dynamic.NewForConfig(kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the metrics client
	metricsclient, err = metrics.NewForConfig(kubeconfig)
	if err != nil {
		panic(err.Error())
	}
}

func PrepareClientset(cfg *rest.Config) {
	// create the clientset
	var err error

	kubeclient, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}

	dynamicinterface, err = dynamic.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}

	// create the metrics client
	metricsclient, err = metrics.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}
}

func InitDiscoveryFromEks(cfg *rest.Config) {
	PrepareClientset(cfg)

	contextbg = context.Background()
}

func InitDiscovery() {
	// var err error
	// kubeclient, err = getClientset()
	// if err != nil {
	// 	panic(err.Error())
	// }

	getClientset()

	contextbg = context.Background()
}

func printJson(category string, jsonMap []map[string]interface{}, isWriteFile bool) {
	jsonObj, err := json.MarshalIndent(jsonMap, "", "\t")
	if err != nil {
		fmt.Printf("ERROR: fail to marshal json(%s), %s\n", category, err.Error())
	}

	//fmt.Printf("[%s]\n%s\n", category, jsonObj)

	if isWriteFile {
		outputpath := "output"
		if _, err := os.Stat(outputpath); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(outputpath, os.ModePerm)
			if err != nil {
				fmt.Println(err)
			}
		}
		os.Remove(outputpath + "/" + strings.ToLower(category) + ".json")
		_ = ioutil.WriteFile(outputpath+"/"+strings.ToLower(category)+".json", jsonObj, 0644)
	}
}

func GetNodesJson() []byte {
	resultMap := make([]map[string]interface{}, 0, 0)

	//namespace := metav1.NamespaceAll
	namespace := metav1.NamespaceDefault
	items, err := GetNodesLikeKubeCtl(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			// https://github.com/kubernetes/client-go/issues/861
			//apiVersion, kind := item.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()
			//singleMap["apiVersion"] = apiVersion
			//singleMap["kind"] = kind
			singleMap["apiVersion"] = "v1"
			singleMap["kind"] = "Node"

			var metadataMap = make(map[string]interface{})
			metadataMap["annotations"] = item.ObjectMeta.Annotations
			metadataMap["creationTimestamp"] = item.ObjectMeta.CreationTimestamp
			metadataMap["generateName"] = item.ObjectMeta.GenerateName
			metadataMap["labels"] = item.ObjectMeta.Labels
			metadataMap["name"] = item.ObjectMeta.Name
			metadataMap["namespace"] = item.ObjectMeta.Namespace
			metadataMap["ownerReferences"] = item.ObjectMeta.OwnerReferences
			metadataMap["resourceVersion"] = item.ObjectMeta.ResourceVersion
			metadataMap["uid"] = item.ObjectMeta.UID
			singleMap["metadata"] = metadataMap

			singleMap["spec"] = item.Spec
			singleMap["status"] = item.Status

			resultMap = append(resultMap, singleMap)
		}
	}

	jsonObj, err := json.MarshalIndent(resultMap, "", "\t")
	if err != nil {
		fmt.Printf("ERROR: fail to marshal json(%s), %s\n", "nodes", err.Error())
		return []byte{}
	}

	return jsonObj
}

func PrintNodesLikeKubectl() {
	resultMap := make([]map[string]interface{}, 0, 0)

	//namespace := metav1.NamespaceAll
	namespace := metav1.NamespaceDefault
	items, err := GetNodesLikeKubeCtl(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			// https://github.com/kubernetes/client-go/issues/861
			//apiVersion, kind := item.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()
			//singleMap["apiVersion"] = apiVersion
			//singleMap["kind"] = kind
			singleMap["apiVersion"] = "v1"
			singleMap["kind"] = "Node"

			var metadataMap = make(map[string]interface{})
			metadataMap["annotations"] = item.ObjectMeta.Annotations
			metadataMap["creationTimestamp"] = item.ObjectMeta.CreationTimestamp
			metadataMap["generateName"] = item.ObjectMeta.GenerateName
			metadataMap["labels"] = item.ObjectMeta.Labels
			metadataMap["name"] = item.ObjectMeta.Name
			metadataMap["namespace"] = item.ObjectMeta.Namespace
			metadataMap["ownerReferences"] = item.ObjectMeta.OwnerReferences
			metadataMap["resourceVersion"] = item.ObjectMeta.ResourceVersion
			metadataMap["uid"] = item.ObjectMeta.UID
			singleMap["metadata"] = metadataMap

			singleMap["spec"] = item.Spec
			singleMap["status"] = item.Status

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("Nodes", resultMap, true)
}

func GetNodesLikeKubeCtl(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]corev1.Node, error) {
	list, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetPodsJson() []byte {
	resultMap := make([]map[string]interface{}, 0, 0)

	//namespace := metav1.NamespaceAll
	namespace := metav1.NamespaceDefault
	items, err := GetPodsLikeKubeCtl(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//apiVersion, kind := item.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()
			//singleMap["apiVersion"] = apiVersion
			//singleMap["kind"] = kind
			singleMap["apiVersion"] = "v1"
			singleMap["kind"] = "Pod"

			var metadataMap = make(map[string]interface{})
			metadataMap["annotations"] = item.ObjectMeta.Annotations
			metadataMap["creationTimestamp"] = item.ObjectMeta.CreationTimestamp
			metadataMap["generateName"] = item.ObjectMeta.GenerateName
			metadataMap["labels"] = item.ObjectMeta.Labels
			metadataMap["name"] = item.ObjectMeta.Name
			metadataMap["namespace"] = item.ObjectMeta.Namespace
			metadataMap["ownerReferences"] = item.ObjectMeta.OwnerReferences
			metadataMap["resourceVersion"] = item.ObjectMeta.ResourceVersion
			metadataMap["uid"] = item.ObjectMeta.UID
			singleMap["metadata"] = metadataMap

			singleMap["spec"] = item.Spec
			singleMap["status"] = item.Status

			resultMap = append(resultMap, singleMap)
		}
	}

	jsonObj, err := json.MarshalIndent(resultMap, "", "\t")
	if err != nil {
		fmt.Printf("ERROR: fail to marshal json(%s), %s\n", "pods", err.Error())
		return []byte{}
	}

	return jsonObj
}

func PrintPodsLikeKubectl() {
	resultMap := make([]map[string]interface{}, 0, 0)

	//namespace := metav1.NamespaceAll
	namespace := metav1.NamespaceDefault
	items, err := GetPodsLikeKubeCtl(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//apiVersion, kind := item.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()
			//singleMap["apiVersion"] = apiVersion
			//singleMap["kind"] = kind
			singleMap["apiVersion"] = "v1"
			singleMap["kind"] = "Pod"

			var metadataMap = make(map[string]interface{})
			metadataMap["annotations"] = item.ObjectMeta.Annotations
			metadataMap["creationTimestamp"] = item.ObjectMeta.CreationTimestamp
			metadataMap["generateName"] = item.ObjectMeta.GenerateName
			metadataMap["labels"] = item.ObjectMeta.Labels
			metadataMap["name"] = item.ObjectMeta.Name
			metadataMap["namespace"] = item.ObjectMeta.Namespace
			metadataMap["ownerReferences"] = item.ObjectMeta.OwnerReferences
			metadataMap["resourceVersion"] = item.ObjectMeta.ResourceVersion
			metadataMap["uid"] = item.ObjectMeta.UID
			singleMap["metadata"] = metadataMap

			singleMap["spec"] = item.Spec
			singleMap["status"] = item.Status

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("Pods", resultMap, true)
}

func GetPodsLikeKubeCtl(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]corev1.Pod, error) {
	list, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// msg := fmt.Sprintf("apiVersion: %s", list.APIVersion)
	// fmt.Println(msg)
	// msg = fmt.Sprintf("kind: %s", list.Kind)
	// fmt.Println(msg)

	return list.Items, nil
}

// type MyMetrics struct {
// 	containerName string `json:"containerName"`
// 	podNamespace string `json:"containerName"`
// 	cpu string `json:"cpu"`
// 	memory string `json:"memory"`
// }

func PrintMetrics() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	// namespace := "cicd-test-k8s-master"

	// Metric 서버가 cluster에 배포된 상태여야 확인 가능
	// https://kubernetes.io/ko/docs/tasks/debug-application-cluster/resource-metrics-pipeline/
	// CPU : 일정 기간 동안 CPU 코어에서 평균 사용량으로 리포트 됨
	// Memory : Metric이 수집된 순간 작업 집합으로 리포트 됨(in-use memory + 일부 cached memory)
	items, err := getMetrics(metricsclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, podMetric := range items {
			//fmt.Printf("%+v\n", podMetric)
			//fmt.Printf("%s\n", podMetric.GetName())
			//fmt.Printf("%s\n", podMetric.GetNamespace())

			for _, container := range podMetric.Containers {
				//fmt.Printf("%vm\n", container.Usage.Cpu().MilliValue())
				//fmt.Printf("%vMi\n", container.Usage.Memory().Value()/(1024*1024))

				cpuQuantity, ok := container.Usage.Cpu().AsInt64()
				memQuantity, ok := container.Usage.Memory().AsInt64()
				if !ok {
					return
				}
				//msg := fmt.Sprintf("Container Name: %s(%s) \n CPU usage: %d \n Memory usage: %d",
				//					container.Name, podMetric.GetNamespace(), cpuQuantity, memQuantity)
				//fmt.Println(msg)

				var singleMap = make(map[string]interface{})
				singleMap["containerName"] = container.Name
				singleMap["podNamespace"] = podMetric.GetNamespace()
				singleMap["cpu"] = strconv.FormatInt(cpuQuantity, 10)
				singleMap["memory"] = strconv.FormatInt(memQuantity, 10)

				resultMap = append(resultMap, singleMap)

				//fmt.Printf("INFO: subResultMap %s\n", subResultMap)
				//fmt.Printf("INFO: resultMap %s\n", resultMap)
			}
		}
	}

	//log.Println("INFO: resultMap %s", resultMap)
	//fmt.Printf("INFO: resultMap %s\n", resultMap)

	//f := colorjson.NewFormatter()
	//f.Indent = 4

	printJson("Metrics", resultMap, true)
}

func getMetrics(clientset *metrics.Clientset, ctx context.Context, namespace string) ([]v1beta1.PodMetrics, error) {
	list, err := clientset.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintNamespace() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetNamespace(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//singleMap["name"] = item.Name
			//singleMap["objectKind"] = item.ObjectKind
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			//singleMap["namespace"] = item.Namespace
			singleMap["uid"] = item.UID
			singleMap["status"] = item.Status

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("Namespace", resultMap, true)
}

func GetNamespace(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]corev1.Namespace, error) {
	list, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func GetServiceAccountJson() []byte {
	resultMap := make([]map[string]interface{}, 0, 0)

	//namespace := metav1.NamespaceAll
	namespace := metav1.NamespaceDefault
	items, err := GetServiceAccount(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			// https://github.com/kubernetes/client-go/issues/861
			//apiVersion, kind := item.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()
			//singleMap["apiVersion"] = apiVersion
			//singleMap["kind"] = kind
			singleMap["apiVersion"] = "v1"
			singleMap["kind"] = "ServiceAccount"

			var metadataMap = make(map[string]interface{})
			metadataMap["creationTimestamp"] = item.ObjectMeta.CreationTimestamp
			metadataMap["name"] = item.ObjectMeta.Name
			metadataMap["namespace"] = item.ObjectMeta.Namespace
			metadataMap["resourceVersion"] = item.ObjectMeta.ResourceVersion
			metadataMap["uid"] = item.ObjectMeta.UID
			singleMap["metadata"] = metadataMap

			resultMap = append(resultMap, singleMap)
		}
	}

	jsonObj, err := json.MarshalIndent(resultMap, "", "\t")
	if err != nil {
		fmt.Printf("ERROR: fail to marshal json(%s), %s\n", "nodes", err.Error())
		return []byte{}
	}

	return jsonObj
}

func PrintServiceAccount() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetServiceAccount(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//singleMap["apiVersion"] = item.APIVersion
			singleMap["kind"] = item.Kind
			singleMap["creationTimestamp"] = item.CreationTimestamp
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace
			singleMap["resourceVersion"] = item.ResourceVersion
			singleMap["uid"] = item.UID
			singleMap["ownerReferences"] = item.OwnerReferences

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("ServiceAccount", resultMap, true)
}

func GetServiceAccount(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]corev1.ServiceAccount, error) {
	list, err := clientset.CoreV1().ServiceAccounts(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintRole() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetRole(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//singleMap["name"] = item.Name
			//singleMap["objectKind"] = item.ObjectKind
			singleMap["kind"] = item.Kind
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace
			singleMap["uid"] = item.UID

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("Role", resultMap, true)
}

func GetRole(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]rbacv1.Role, error) {
	list, err := clientset.RbacV1().Roles(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintClusterRole() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetClusterRole(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//singleMap["name"] = item.Name
			//singleMap["objectKind"] = item.ObjectKind
			singleMap["kind"] = item.Kind
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace
			singleMap["uid"] = item.UID

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("ClusterRole", resultMap, true)
}

func GetClusterRole(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]rbacv1.ClusterRole, error) {
	list, err := clientset.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintRoleBinding() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetRoleBinding(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//singleMap["name"] = item.Name
			//singleMap["objectKind"] = item.ObjectKind
			singleMap["kind"] = item.Kind
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace
			singleMap["uid"] = item.UID

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("RoleBinding", resultMap, true)
}

func GetRoleBinding(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]rbacv1.RoleBinding, error) {
	list, err := clientset.RbacV1().RoleBindings(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintClusterRoleBinding() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetClusterRoleBinding(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//singleMap["name"] = item.Name
			//singleMap["objectKind"] = item.ObjectKind
			singleMap["kind"] = item.Kind
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace
			singleMap["uid"] = item.UID

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("ClusterRoleBinding", resultMap, true)
}

func GetClusterRoleBinding(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]rbacv1.ClusterRoleBinding, error) {
	list, err := clientset.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintPods() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetPods(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			// containers := item.Spec.Containers
			// for _, container := range containers {
			// 	fmt.Printf("ContainerName: %s\nContainerImage: %s\n", container.Name, container.Image)
			// }

			if len(item.OwnerReferences) == 0 {
				fmt.Printf("Pod %s has no owner", item.Name)
				continue
			}

			var ownerName, ownerKind string

			switch item.OwnerReferences[0].Kind {
			case "ReplicaSet":
				replica, repErr := kubeclient.AppsV1().ReplicaSets(item.Namespace).Get(context.TODO(), item.OwnerReferences[0].Name, metav1.GetOptions{})
				if repErr != nil {
					panic(repErr.Error())
				}

				if replica.OwnerReferences != nil {
					ownerName = replica.OwnerReferences[0].Name
					ownerKind = "Deployment"
				} else {
					// exception
					fmt.Println("replica.OwnerReferences is nil")
				}
			case "DaemonSet", "StatefulSet":
				ownerName = item.OwnerReferences[0].Name
				ownerKind = item.OwnerReferences[0].Kind
			//case "Node":
			//	ownerName = item.OwnerReferences[0].Name
			//	ownerKind = item.OwnerReferences[0].Kind
			default:
				fmt.Printf("Could not find resource manager for type %s\n", item.OwnerReferences[0].Kind)
				//fmt.Printf("Name(%s)\n", item.OwnerReferences[0].Name)
				continue
			}

			//fmt.Printf("POD %s is managed by %s %s\n", item.Name, ownerName, ownerKind)
			var singleMap = make(map[string]interface{})
			singleMap["name"] = item.Name
			singleMap["ownerName"] = ownerName
			singleMap["ownerKind"] = ownerKind

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("Pods", resultMap, true)
}

func GetPods(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]corev1.Pod, error) {
	list, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintContainers() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetPods(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			containers := item.Spec.Containers
			for _, container := range containers {
				//fmt.Printf("ContainerName: %s\nContainerImage: %s\n", container.Name, container.Image)

				var singleMap = make(map[string]interface{})
				singleMap["containerName"] = container.Name
				singleMap["containerImage"] = container.Image
				singleMap["podName"] = item.Name

				resultMap = append(resultMap, singleMap)
			}
		}
	}

	printJson("Containers", resultMap, true)
}

func PrintDeployments() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetDeployments(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			singleMap["kind"] = item.Kind
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("Deployments", resultMap, true)
}

func GetDeployments(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
	list, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintResources() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	//items, err := GetResources(kubeclient, contextbg, "apps", "v1", "deployments", namespace)
	items, err := GetResourcesDynamically(dynamicinterface, contextbg, "apps", "v1", "deployments", namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			singleMap["apiVersion"] = item.GetAPIVersion()
			singleMap["kind"] = item.GetKind()
			singleMap["labels"] = item.GetLabels()
			singleMap["name"] = item.GetName()
			singleMap["namespace"] = item.GetNamespace()
			singleMap["uid"] = item.GetUID()

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("Resources", resultMap, true)
}

/*
func GetResources(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
	list, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
*/

func GetResourcesDynamically(dynamic dynamic.Interface, ctx context.Context,
	group string, version string, resource string, namespace string) ([]unstructured.Unstructured, error) {

	resourceId := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	list, err := dynamic.Resource(resourceId).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

/*
func PrintResourcesByJq() {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// config를 아래 함수로 구하면 runtime에 kubeconfig redefined 에러가 발생
	//config := ctrl.GetConfigOrDie()

	// create the kubeClient
	//clientset, err := kubernetes.NewForConfig(config)
	//dynamic, err := dynamic.NewForConfig(config)
	//if err != nil {
	//	panic(err.Error())
	//}
	dynamic := dynamic.NewForConfigOrDie(config)

	ctx := context.Background()
	namespace := metav1.NamespaceAll
	//query := ".metadata.labels[\"app.kubernetes.io/managed-by\"] == \"Helm\""
	query := "TIGERA_OPERATOR_INIT_IMAGE_VERSION"
	items, err := GetResourcesByJq(dynamic, ctx, "apps", "v1", "deployments", namespace, query)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			fmt.Printf("%+v\n", item)
		}
	}
}

func GetResourcesByJq(dynamic dynamic.Interface, ctx context.Context,
	group string, version string, resource string, namespace string, jq string) ([]unstructured.Unstructured, error) {

	resources := make([]unstructured.Unstructured, 0)
	query, err := gojq.Parse(jq)
	if err != nil {
		return nil, err
	}

	items, err := GetResourcesDynamically(dynamic, ctx, group, version, resource, namespace)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		// Convert object to raw JSON
		var rawJson interface{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, &rawJson)
		if err != nil {
			return nil, err
		}

		//fmt.Println(rawJson)
		// Evaluate jq against JSON
		iter := query.Run(rawJson)
		for {
			result, ok := iter.Next()
			if !ok {
				break
			}

			if err, ok := result.(error); ok {
				if err != nil {
					return nil, err
				}
			} else {
				boolResult, ok := result.(bool)
				if !ok {
					fmt.Println("Query returned non-boolean value")
				} else if boolResult {
					resources = append(resources, item)
				}
			}
		}
	}

	return resources, nil
}
*/

func PrintServices() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetServices(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//singleMap["objectKind"] = item.GetObjectKind()
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace
			singleMap["uid"] = item.UID
			singleMap["spec"] = item.Spec
			singleMap["status"] = item.Status

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("Services", resultMap, true)
}

func GetServices(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]corev1.Service, error) {
	list, err := clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintReplicaSets() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetReplicaSets(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//singleMap["objectKind"] = item.GetObjectKind()
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace
			singleMap["uid"] = item.UID

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("ReplicaSets", resultMap, true)
}

func GetReplicaSets(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]appsv1.ReplicaSet, error) {
	list, err := clientset.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintDaemonSets() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetDaemonSets(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//singleMap["objectKind"] = item.GetObjectKind()
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace
			singleMap["uid"] = item.UID
			singleMap["spec"] = item.Spec
			singleMap["status"] = item.Status

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("DaemonSets", resultMap, true)
}

func GetDaemonSets(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]appsv1.DaemonSet, error) {
	list, err := clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintStatefulSets() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetStatefulSets(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//singleMap["objectKind"] = item.GetObjectKind()
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace
			singleMap["uid"] = item.UID
			singleMap["spec"] = item.Spec
			singleMap["status"] = item.Status

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("StatefulSets", resultMap, true)
}

func GetStatefulSets(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]appsv1.StatefulSet, error) {
	list, err := clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintConfigMaps() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetConfigMaps(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//singleMap["objectKind"] = item.GetObjectKind()
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace
			singleMap["uid"] = item.UID

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("ConfigMaps", resultMap, true)
}

func GetConfigMaps(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]corev1.ConfigMap, error) {
	list, err := clientset.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintSecrets() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetSecrets(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//singleMap["objectKind"] = item.GetObjectKind()
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace
			singleMap["uid"] = item.UID

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("Secrets", resultMap, true)
}

func GetSecrets(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]corev1.Secret, error) {
	list, err := clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintPersistentVolume() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetPersistentVolume(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			singleMap["apiVersion"] = item.APIVersion
			singleMap["creationTimestamp"] = item.CreationTimestamp
			//singleMap["objectKind"] = item.GetObjectKind()
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace
			singleMap["uid"] = item.UID

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("PersistentVolume", resultMap, true)
}

func GetPersistentVolume(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]corev1.PersistentVolume, error) {
	list, err := clientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintPersistentVolumeClaim() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetPersistentVolumeClaim(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			singleMap["apiVersion"] = item.APIVersion
			singleMap["creationTimestamp"] = item.CreationTimestamp
			//singleMap["objectKind"] = item.GetObjectKind()
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace
			singleMap["uid"] = item.UID

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("PersistentVolumeClaim", resultMap, true)
}

func GetPersistentVolumeClaim(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]corev1.PersistentVolumeClaim, error) {
	list, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintIngress() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetIngress(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//singleMap["objectKind"] = item.GetObjectKind()
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace
			singleMap["uid"] = item.UID

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("Ingress", resultMap, true)
}

func GetIngress(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]networkv1.Ingress, error) {
	list, err := clientset.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func PrintNetworkPolicy() {
	resultMap := make([]map[string]interface{}, 0, 0)

	namespace := metav1.NamespaceAll
	items, err := GetNetworkPolicy(kubeclient, contextbg, namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range items {
			//fmt.Printf("%+v\n", item)

			var singleMap = make(map[string]interface{})
			//singleMap["objectKind"] = item.GetObjectKind()
			singleMap["labels"] = item.Labels
			singleMap["name"] = item.Name
			singleMap["namespace"] = item.Namespace
			singleMap["uid"] = item.UID
			singleMap["podSelector"] = item.Spec.PodSelector

			resultMap = append(resultMap, singleMap)
		}
	}

	printJson("NetworkPolicy", resultMap, true)
}

func GetNetworkPolicy(clientset *kubernetes.Clientset, ctx context.Context, namespace string) ([]networkv1.NetworkPolicy, error) {
	list, err := clientset.NetworkingV1().NetworkPolicies(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
