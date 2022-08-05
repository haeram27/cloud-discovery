package k8sapis

import (
	aws "awsdisc/client/aws/apis"
	"testing"
)

func TestDiscoverK8S(t *testing.T) {
	cluster, err := aws.EKSDescribeClusterCmd(aws.AwsConfig(), "eks-cicd-sec-test-ec2-ssh")
	if err != nil {
		t.Error(err)
	}

	k8scfg, err := aws.EKSK8sConfig(cluster.Cluster)
	if err != nil {
		t.Error(err)
	}

	InitDiscoveryFromEks(k8scfg)
	t.Log(string(GetServiceAccountJson()))
	t.Log(string(GetNodesJson()))
	t.Log(string(GetPodsJson()))
}
