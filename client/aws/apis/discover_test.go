package apis

import (
	"awsdisc/apps/util"
	"testing"
)

func TestDiscoverAll(t *testing.T) {
	j, err := DiscoverAll()
	if err != nil {
		t.Error(err)
	}
	t.Log(util.PrettyJson(j).String())
}
