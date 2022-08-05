package apis

import (
	"context"
	"encoding/json"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	resourcemanagerpb "google.golang.org/genproto/googleapis/cloud/resourcemanager/v3"
)

// List Projects from folders or organizations
// parent: folders/{name} or organizations/{name}
func GCPAPIResourceMngrListProjects(parent string) ([]string, error) {
	ctx := context.Background()
	c, err := resourcemanager.NewProjectsClient(ctx, option.WithCredentialsFile("/home/swvm/.config/gcloud/credentials/poetic-diorama-358105-8d9d41114576.json"))
	if err != nil {
		return []string{}, err
	}
	defer c.Close()

	// folders/{id}, organizations/{id}
	req := &resourcemanagerpb.ListProjectsRequest{
		// TODO: Fill request struct fields.
		// See https://cloud.google.com/resource-manager/reference/rest/v3/projects/list
		// See https://pkg.go.dev/google.golang.org/genproto/googleapis/cloud/resourcemanager/v3#ListProjectsRequest.
		Parent: parent,
	}

	var ret []string
	if it := c.ListProjects(ctx, req); it != nil {
		for {
			project, err := it.Next()
			if err == iterator.Done {
				break
			}

			if err != nil {
				return []string{}, err
			}

			if blob, err := json.Marshal(project); err != nil {
				return []string{}, err
			} else {
				ret = append(ret, string(blob))
			}
		}
	}

	return ret, nil
}

// Search Active Projects
func GCPAPIResourceMngrSearchProjects() ([]string, error) {
	ctx := context.Background()
	c, err := resourcemanager.NewProjectsClient(ctx, option.WithCredentialsFile("/home/swvm/.config/gcloud/credentials/poetic-diorama-358105-8d9d41114576.json"))
	if err != nil {
		return []string{}, err
	}
	defer c.Close()

	// folders/{id}, organizations/{id}
	req := &resourcemanagerpb.SearchProjectsRequest{}

	var ret []string
	if it := c.SearchProjects(ctx, req); it != nil {
		for {
			project, err := it.Next()
			if err == iterator.Done {
				break
			}

			if err != nil {
				return []string{}, err
			}

			if blob, err := json.Marshal(project); err != nil {
				return []string{}, err
			} else {
				ret = append(ret, string(blob))
			}
		}
	}

	return ret, nil
}
