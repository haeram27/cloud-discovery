package apis

import (
	"testing"
)

func TestGCPResourceMngrListProjects(t *testing.T) {
	// parent: folders/{name}  or organizations/{name}
	if arr, err := GCPAPIResourceMngrListProjects("organizations/447611139380"); err != nil {
		t.Error(err)
	} else {
		t.Log(arr)
	}
}

func TestGCPResourceMngrSearchProjects(t *testing.T) {
	if arr, err := GCPAPIResourceMngrSearchProjects(); err != nil {
		t.Error(err)
	} else {
		t.Log(arr)
	}
}
