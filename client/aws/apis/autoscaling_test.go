package apis

import (
	"awsdisc/apps/util"
	"encoding/json"
	"testing"
)

func TestASCLDescribeAutoScalingGroupsAll(t *testing.T) {
	var jsonBlob []byte

	result, err := ASCLDescribeAutoScalingGroupsCmd(AwsConfig())
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	//t.Log(util.PrettyJson(jsonBlob).String())

	values := util.JsonPath(jsonBlob, "$.AutoScalingGroups[*].Instances[0].InstanceId")
	//value, err := util.JsonPath(jsonBlob, "$.AutoScalingGroups[*].Instances[1:].InstanceId"))

	t.Logf("%v", values)

}

func TestJsonObjectPath(t *testing.T) {
	rawjson := `{ "key" : "value" }`

	values := util.JsonPath([]byte(rawjson), "$.key")

	t.Logf("%+v", values)
}
