package apis

import (
	"cloudisc/apps/util"
	"encoding/json"
	"testing"
)

func TestEksDiscoveryAll(t *testing.T) {
	var jsonBlob []byte
	var result interface{}

	result, err := EKSListClustersCmd(AwsConfig())
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())

	result, err = EKSDescribeClusterCmd(AwsConfig(), "eks-cicd-sec-test-ec2-ssh")
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())

	result, err = EKSDescribeClusterCmd(AwsConfig(), "private-subnet-cluster")
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())

	/*
		j, err = EKSDescribeClusterCmd(AwsConfig(), "")
		if err != nil {
			t.Error(err)
		}
		t.Log(util.PrettyJson(j).String())
	*/
}

func TestEksConfig(t *testing.T) {
	result, err := EKSDescribeClusterCmd(AwsConfig(), "eks-cicd-sec-test-ec2-ssh")
	if err != nil {
		t.Error(err)
	}
	EKSK8sClientset(result.Cluster)
}
