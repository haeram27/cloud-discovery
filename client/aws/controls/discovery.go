package controls

import (
	"cloudisc/apps"
	"cloudisc/apps/util"
	aws "cloudisc/client/aws/apis"
	"encoding/json"
	"fmt"
)

func ListEc2InstancesForEBSScan() []string {
	r1, err := aws.EC2DescribeInstancesCmd(aws.AwsConfig())
	if err != nil {
		apps.Logs.Error(err)
		return []string{}
	}

	r2, err := aws.ASCLDescribeAutoScalingGroupsCmd(aws.AwsConfig())
	if err != nil {
		apps.Logs.Error(err)
		return []string{}
	}

	instances, err := json.Marshal(r1)
	if err != nil {
		apps.Logs.Error(err)
		return []string{}
	}

	groups, err := json.Marshal(r2)
	if err != nil {
		apps.Logs.Error(err)
		return []string{}
	}

	j1 := util.JsonPath([]byte(instances), "$.Reservations[*].Instances[*].InstanceId")
	if err != nil {
		apps.Logs.Error(err)
		return []string{}
	}

	j2 := util.JsonPath([]byte(groups), "$.AutoScalingGroups[*].Instances[1:].InstanceId")
	if err != nil {
		apps.Logs.Error(err)
		return []string{}
	}

	s := make([]string, len(j1))
	for i, v := range j1 {
		s[i] = fmt.Sprint(v)
	}

	for _, e := range j2 {
		s = util.RemoveFromSlice(s, e.(string))
	}

	return s
}
