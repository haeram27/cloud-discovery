package apis

import (
	"cloudisc/apps"
	"cloudisc/apps/util"
	"testing"
)

func TestGCPComputeListInstances(t *testing.T) {
	instances, err := GCPAPIComputeListInstances("", "")
	if err != nil {
		t.Error(err)
	}

	for _, instance := range instances {
		apps.Logs.Debug(util.PrettyJson([]byte(instance)))
	}

}

func TestGCPAPIComputeAggregatedListInstances(t *testing.T) {
	instances, err := GCPAPIComputeAggregatedListInstances("poetic-diorama-358105")
	if err != nil {
		t.Error(err)
	}

	for _, instance := range instances {
		apps.Logs.Debug(util.PrettyJson([]byte(instance)))
	}
}
