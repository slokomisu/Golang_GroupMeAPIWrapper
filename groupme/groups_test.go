package groupme

import (
	"testing"
	"os"
)

func TestGroupsService_Index(t *testing.T) {
	service := &groupsService{
		token: os.Getenv("ACCESS_TOKEN"),
		baseUrl: "https://api.groupme.com/v3",
	}
	groups, err := service.Index(nil)

	if err != nil {
		t.Failed()
	}

	if groups == nil {
		t.Failed()
	} else if len(groups) != 10 {
		t.Failed()
	}

}

func TestGroupsService_Show(t *testing.T) {
	service := &groupsService{
		token: os.Getenv("ACCESS_TOKEN"),
		baseUrl: "https://api.groupme.com/v3",
	}
	group, err := service.Show(32968213)
	if err != nil {
		t.Failed()
	}
	if group == nil {
		t.Failed()
	} else if group.Name != "Test Group" {
		t.Failed()
	}

}
