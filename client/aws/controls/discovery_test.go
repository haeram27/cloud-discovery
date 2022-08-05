package controls

import (
	"awsdisc/apps/util"
	"fmt"
	"testing"
)

func TestScanTarget(t *testing.T) {
	s := ListEc2InstancesForEBSScan()
	t.Log(s)
}

func TestScanTargetExam(t *testing.T) {
	t1 := `{
		"Reservations": [
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

	t2 := `{
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

	j1 := util.JsonPath([]byte(t1), "$.Reservations[*].Instances[*].InstanceId")
	t.Log(j1)

	j2 := util.JsonPath([]byte(t2), "$.AutoScalingGroups[*].Instances[1:].InstanceId")
	t.Log(j2)

	s := make([]string, len(j1))
	for i, v := range j1 {
		s[i] = fmt.Sprint(v)
	}

	for _, e := range j2 {
		s = util.RemoveFromSlice(s, e.(string))
	}

	t.Log(s)
}
