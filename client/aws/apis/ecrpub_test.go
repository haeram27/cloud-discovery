package apis

import (
	"awsdisc/apps/util"
	"encoding/json"
	"testing"
)

func TestEcrPubDiscoverAll(t *testing.T) {
	var jsonBlob []byte
	var result interface{}

	result, err := ECRPubDescribeRegistryCmd(AwsConfig())
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
