package apis

import (
	"context"
	"encoding/json"
	"errors"

	compute "cloud.google.com/go/compute/apiv1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

// listInstances prints a list of instances created in given project in given zone.
// $ gcloud compute instances list
func GCPAPIComputeListInstances(projectID, zone string) ([]string, error) {
	if projectID == "" {
		// projectID = "poetic-diorama-358105"
		return []string{}, errors.New("invalid arguments: projectid")
	}

	if zone == "" {
		// zone = "asia-northeast3-a"
		return []string{}, errors.New("invalid arguments: zone")
	}

	ctx := context.Background()
	c, err := compute.NewInstancesRESTClient(ctx, option.WithCredentialsFile("/home/swvm/.config/gcloud/credentials/poetic-diorama-358105-8d9d41114576.json"))
	if err != nil {
		return []string{}, err
	}
	defer c.Close()

	req := &computepb.ListInstancesRequest{
		Project: projectID,
		Zone:    zone,
	}

	var ret []string
	if it := c.List(ctx, req); it != nil {
		for {
			instance, err := it.Next()
			if err == iterator.Done {
				break
			}

			if err != nil {
				return []string{}, err
			}

			if blob, err := json.Marshal(instance); err != nil {
				return []string{}, err
			} else {
				ret = append(ret, string(blob))
			}
		}
	}

	return ret, nil
}

func GCPAPIComputeAggregatedListInstances(projectID string) ([]string, error) {

	ctx := context.Background()
	c, err := compute.NewInstancesRESTClient(ctx, option.WithCredentialsFile("/home/swvm/.config/gcloud/credentials/poetic-diorama-358105-8d9d41114576.json"))
	if err != nil {
		return []string{}, err
	}
	defer c.Close()

	req := &computepb.AggregatedListInstancesRequest{
		Project: projectID,
	}

	var ret []string
	if it := c.AggregatedList(ctx, req); it != nil {
		for {
			pair, err := it.Next()
			if err == iterator.Done {
				break
			}

			if err != nil {
				return []string{}, err
			}

			if pair.Value.Instances != nil {
				for _, instance := range pair.Value.Instances {
					if blob, err := json.Marshal(instance); err != nil {
						return []string{}, err
					} else {
						ret = append(ret, string(blob))
					}
				}
			}
		}
	}

	return ret, nil
}
