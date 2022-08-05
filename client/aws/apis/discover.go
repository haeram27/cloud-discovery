package apis

import (
	apps "awsdisc/apps"
	"awsdisc/client/k8sapis"
	"encoding/json"
)

func DiscoverAll() ([]byte, error) {
	var result interface{}
	var ec2, ecr, ecs, eks, k8snodes, k8spods []byte

	// EC2
	result, err := EC2DescribeInstancesCmd(AwsConfig())
	if err != nil {
		apps.Logs.Error(err)
	} else {
		ec2, err = json.Marshal(result)
		if err != nil {
			apps.Logs.Error(err)
		}
	}

	// ECR
	result, err = ECRDescribeRepositoriesCmd(AwsConfig())
	if err != nil {
		apps.Logs.Error(err)
	} else {
		ecr, err = json.Marshal(result)
		if err != nil {
			apps.Logs.Error(err)
		}
	}

	// ECS
	result, err = ECSDescribeClustersCmd(AwsConfig(), []string{"cicd-ecs-ec2-cluster", "swh-ecs-cluster-ssh", "cicd-ecs-cluster"})
	if err != nil {
		apps.Logs.Error(err)
	} else {
		ecs, err = json.Marshal(result)
		if err != nil {
			apps.Logs.Error(err)
		}
	}

	// EKS
	result, err = EKSListClustersCmd(AwsConfig())
	if err != nil {
		apps.Logs.Error(err)
	} else {
		eks, err = json.Marshal(result)
		if err != nil {
			apps.Logs.Error(err)
		}
	}

	cluster, err := EKSDescribeClusterCmd(AwsConfig(), "eks-cicd-sec-test-ec2-ssh")
	if err != nil {
		apps.Logs.Error(err)
	}

	k8scfg, err := EKSK8sConfig(cluster.Cluster)
	if err != nil {
		apps.Logs.Error(err)
	}

	if k8scfg != nil {
		k8sapis.InitDiscoveryFromEks(k8scfg)
		k8snodes = k8sapis.GetNodesJson()
		k8spods = k8sapis.GetPodsJson()
	}

	// consolidate
	data := make(map[string]interface{})
	var inner interface{}
	json.Unmarshal(ec2, &inner)
	data["EC2"] = inner
	json.Unmarshal(ecr, &inner)
	data["ECRRepositories"] = inner
	json.Unmarshal(ecs, &inner)
	data["ECSClusters"] = inner
	json.Unmarshal(eks, &inner)
	data["EKSClusters"] = inner
	json.Unmarshal(k8snodes, &inner)
	data["EKSNodes"] = inner
	json.Unmarshal(k8spods, &inner)
	data["EKSPods"] = inner

	discovered, err := json.Marshal(data)
	if err != nil {
		apps.Logs.Error(err)
		return nil, err
	}

	return discovered, nil
}
