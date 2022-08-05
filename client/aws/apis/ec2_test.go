package apis

import (
	"awsdisc/apps/util"
	"encoding/json"
	"testing"
)

func TestEC2DiscoverAll(t *testing.T) {
	var jsonBlob []byte
	var result interface{}

	result, err := EC2DescribeInstancesCmd(AwsConfig())
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

func TestEC2CreateSnapshots(t *testing.T) {
	var jsonBlob []byte
	var result interface{}

	result, err := EC2CreateSnapshotsCmd(AwsConfig(), "i-0781a799411110be1")
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
