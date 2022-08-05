package apis

import (
	apps "cloudisc/apps"
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func EC2DescribeInstancesCmd(cfg *aws.Config) (*ec2.DescribeInstancesOutput, error) {
	if cfg == nil || cfg.Credentials == nil {
		err := errors.New("invalid aws config")
		apps.Logs.Error(err)
		return nil, err
	}

	client := ec2.NewFromConfig(*cfg)
	if client == nil {
		err := errors.New("failed to initialize aws client")
		apps.Logs.Error(err)
		return nil, err
	}

	awsctx := context.TODO()
	input := &ec2.DescribeInstancesInput{}

	return client.DescribeInstances(awsctx, input)
}

func EC2CreateSnapshotsCmd(cfg *aws.Config, instanceId string) (*ec2.CreateSnapshotsOutput, error) {
	if cfg == nil || cfg.Credentials == nil {
		err := errors.New("invalid aws config")
		apps.Logs.Error(err)
		return nil, err
	}

	client := ec2.NewFromConfig(*cfg)
	if client == nil {
		err := errors.New("failed to initialize aws client")
		apps.Logs.Error(err)
		return nil, err
	}

	ispc := types.InstanceSpecification{
		InstanceId: &instanceId,
	}

	awsctx := context.TODO()
	description := "cicd-sec"
	tagKey := "instanceId"
	input := &ec2.CreateSnapshotsInput{
		InstanceSpecification: &ispc,
		Description:           &description,
		TagSpecifications:     []types.TagSpecification{{ResourceType: "snapshot", Tags: []types.Tag{{Key: &tagKey, Value: &instanceId}}}},
	}

	return client.CreateSnapshots(awsctx, input)
}
