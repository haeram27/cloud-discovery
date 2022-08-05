package controls

import (
	"awsdisc/apps"
	"awsdisc/apps/util"
	"awsdisc/client/gcp/apis"
)

func GCPListInstances() error {
	projects, err := apis.GCPAPIResourceMngrSearchProjects()
	if err != nil {
		return err
	}

	for _, project := range projects {
		for _, id := range util.JsonPath([]byte(project), "$.project_id") {
			instances, err := apis.GCPAPIComputeAggregatedListInstances(id.(string))
			if err != nil {
				return err
			}
			for _, instance := range instances {
				apps.Logs.Debug(util.PrettyJson([]byte(instance)))
			}
		}
	}

	return nil
}
