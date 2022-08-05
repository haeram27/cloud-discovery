package apis

import (
	"cloudisc/apps/util"
	"encoding/json"
	"testing"
)

func TestEcsDiscoverAll(t *testing.T) {
	var jsonBlob []byte
	var result interface{}

	result, err := ECSListClustersCmd(AwsConfig())
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())

	result, err = ECSDescribeClustersCmd(AwsConfig(), []string{"cicd-ecs-ec2-cluster", "swh-ecs-cluster-ssh", "cicd-ecs-cluster"})
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())

	result, err = ECSListContainerInstancesCmd(AwsConfig(), "cicd-ecs-ec2-cluster")
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())

	result, err = ECSListContainerInstancesCmd(AwsConfig(), "swh-ecs-cluster-ssh")
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())

	result, err = ECSListContainerInstancesCmd(AwsConfig(), "cicd-ecs-cluster")
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())

	result, err = ECSListTaskDefinitionsCmd(AwsConfig())
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())

	result, err = ECSDescribeTaskDefinitionCmd(AwsConfig(), "cicd-task-nginx:1")
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())

	result, err = ECSDescribeTaskDefinitionCmd(AwsConfig(), "cicd-task-ubuntu_nginx:2")
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())

	result, err = ECSDescribeTaskDefinitionCmd(AwsConfig(), "sw-task:4")
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())
}

func TestECSListClusters(t *testing.T) {
	var jsonBlob []byte
	var result interface{}

	result, err := ECSListClustersCmd(AwsConfig())
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())
}

func TestECSDescribeClusters(t *testing.T) {
	var jsonBlob []byte
	var result interface{}

	result, err := ECSDescribeClustersCmd(AwsConfig(), []string{"cicd-ecs-ec2-cluster", "swh-ecs-cluster-ssh", "cicd-ecs-cluster"})
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())
}

func TestECSListContainerInstances(t *testing.T) {
	var jsonBlob []byte
	var result interface{}

	result, err := ECSListContainerInstancesCmd(AwsConfig(), "cicd-ecs-ec2-cluster")
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())

	result, err = ECSListContainerInstancesCmd(AwsConfig(), "swh-ecs-cluster-ssh")
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())

	result, err = ECSListContainerInstancesCmd(AwsConfig(), "cicd-ecs-cluster")
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())
}
