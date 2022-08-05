package apis

import (
	"awsdisc/apps/util"
	"encoding/json"
	"testing"
)

func TestEcrDiscoverAll(t *testing.T) {
	var jsonBlob []byte
	var result interface{}

	result, err := ECRDescribeRegistryCmd(AwsConfig())
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())

	result, err = ECRDescribeRepositoriesCmd(AwsConfig())
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

func TestECRDescribeRegistryCmd(t *testing.T) {
	var jsonBlob []byte
	var result interface{}

	result, err := ECRDescribeRegistryCmd(AwsConfig())
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

func TestECRDescribeRepositoriesCmd(t *testing.T) {
	var jsonBlob []byte
	var result interface{}

	result, err := ECRDescribeRepositoriesCmd(AwsConfig())
	if err != nil {
		t.Error(err)
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log(util.PrettyJson(jsonBlob).String())

	repos := util.JsonPath(jsonBlob, "$.Repositories[:].RepositoryName")
	t.Logf("%v", repos)
}

func TestECRListImagesAll(t *testing.T) {
	ECRListImagesAll(AwsConfig())
}

func TestECRListImagesAllST(t *testing.T) {
	images := ECRListImagesAllST(AwsConfig())

	for _, image := range images {
		t.Log(EcrImageUri(image).DigestUri())
	}
}

func TestECRGetAuthorizationTokenCmd(t *testing.T) {
	tok, err := ECRGetAuthorizationTokenCmd(AwsConfig())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tok)
}
