package util

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestPrintPrettyJson(t *testing.T) {
	inner1 := `{"Name": "iron man"}`
	inner2 := `{"Type": "developer"}`

	// consolidate
	jsonStrt := make(map[string]interface{})
	var innerJson interface{}
	json.Unmarshal([]byte(inner1), &innerJson)
	jsonStrt["inner1"] = innerJson
	json.Unmarshal([]byte(inner2), &innerJson)
	jsonStrt["inner2"] = innerJson

	jsonBlob, err := json.Marshal(jsonStrt)
	if err != nil {
		t.Error(err)
		return
	}

	PrintPrettyJson(jsonBlob)
}

func TestPrettyJson(t *testing.T) {
	jsonBlob := `{
		Name: "iron man",
		Type: "developer",
	}`

	t.Log(PrettyJson([]byte(jsonBlob)).String())
}

func TestPrettyMarshalJson(t *testing.T) {
	jsonStrt := struct {
		Name string
		Type string
	}{
		Name: "iron man",
		Type: "developer",
	}

	jsonBlob, err := json.MarshalIndent(jsonStrt, "", "  ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(jsonBlob))
}

func TestJsonPath(t *testing.T) {
	jsonBlob := `{ 
		"Name": { 
			"First": "man",
			"Last": "iron"
		},
		"Type": "developer"
	}`

	path := "$.Name.Last"

	v := JsonPath([]byte(jsonBlob), path)
	t.Logf("%+v", v)
}

func TestJsonPathDiff(t *testing.T) {
	rawjson := `{
		"AutoScalingGroups": [
			{
				"Instances": [
					{
						"InstanceId": "i-11"
					},
					{
						"InstanceId": "i-12"
					},
					{
						"InstanceId": "i-13"
					}
				]
			},
			{
				"Instances": [
					{
						"InstanceId": "i-21"
					},
					{
						"InstanceId": "i-22"
					},
					{
						"InstanceId": "i-23"
					}
				]
			},
			{
				"Instances": [
					{
						"InstanceId": "i-31"
					},
					{
						"InstanceId": "i-32"
					},
					{
						"InstanceId": "i-33"
					}
				]
			}
		]
	}`

	rawjson2 := `{
		"AutoScalingGroups": [
			{
				"Instances": [
					{
						"InstanceId": "i-12"
					},
					{
						"InstanceId": "i-13"
					}
				]
			},
			{
				"Instances": [
					{
						"InstanceId": "i-22"
					},
					{
						"InstanceId": "i-23"
					}
				]
			},
			{
				"Instances": [
					{
						"InstanceId": "i-32"
					},
					{
						"InstanceId": "i-33"
					}
				]
			}
		]
	}`

	rawjson3 := `[]`

	values := JsonPath([]byte(rawjson), "$.AutoScalingGroups[*].Instances[*].InstanceId")
	values2 := JsonPath([]byte(rawjson2), "$.AutoScalingGroups[*].Instances[*].InstanceId")
	values3 := JsonPath([]byte(rawjson3), "$.AutoScalingGroups[*].Instances[*].InstanceId")

	_ = rawjson
	_ = rawjson2
	_ = rawjson3
	_ = values3

	start := time.Now()
	s := make([]string, len(values))
	for i, v := range values {
		s[i] = fmt.Sprint(v)
	}

	// JsonPath's return type is []interface or []interface{}
	for _, e := range values2 {
		s = RemoveFromSlice(s, e.(string))
	}

	elipsed := time.Since(start)
	t.Log(elipsed, "=============================================================")
	t.Log(s)

	start = time.Now()
	r1 := make([]string, len(values))
	for i, v := range values {
		r1[i] = fmt.Sprint(v)
	}

	r2 := make([]string, len(values2))
	for i, v := range values2 {
		r2[i] = fmt.Sprint(v)
	}

	r3 := DiffStrSlice(r1, r2)

	elipsed = time.Since(start)
	t.Log(elipsed, "=============================================================")
	t.Log(r3)
}
