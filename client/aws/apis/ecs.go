package apis

import (
	apps "awsdisc/apps"
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func ECSDescribeClustersCmd(cfg *aws.Config, name []string) (*ecs.DescribeClustersOutput, error) {
	if cfg == nil || cfg.Credentials == nil {
		err := errors.New("invalid aws config")
		apps.Logs.Error(err)
		return nil, err
	}

	if len(name) == 0 {
		err := errors.New("invalid arguments: empty name")
		apps.Logs.Error(err)
		return nil, err
	}

	client := ecs.NewFromConfig(*cfg)
	if client == nil {
		err := errors.New("failed to initialize aws client")
		apps.Logs.Error(err)
		return nil, err
	}

	awsctx := context.TODO()
	input := &ecs.DescribeClustersInput{}
	input.Clusters = name
	return client.DescribeClusters(awsctx, input)
}

func ECSDescribeTaskDefinitionCmd(cfg *aws.Config, task string) (*ecs.DescribeTaskDefinitionOutput, error) {
	if cfg == nil || cfg.Credentials == nil {
		err := errors.New("invalid aws config")
		apps.Logs.Error(err)
		return nil, err
	}

	if task == "" {
		err := errors.New("invalid arguments: empty name")
		apps.Logs.Error(err)
		return nil, err
	}

	client := ecs.NewFromConfig(*cfg)
	if client == nil {
		err := errors.New("failed to initialize aws client")
		apps.Logs.Error(err)
		return nil, err
	}

	awsctx := context.TODO()
	input := &ecs.DescribeTaskDefinitionInput{}
	input.TaskDefinition = &task
	return client.DescribeTaskDefinition(awsctx, input)
}

func ECSListClustersCmd(cfg *aws.Config) (*ecs.ListClustersOutput, error) {
	if cfg == nil || cfg.Credentials == nil {
		err := errors.New("invalid aws config")
		apps.Logs.Error(err)
		return nil, err
	}

	client := ecs.NewFromConfig(*cfg)
	if client == nil {
		err := errors.New("failed to initialize aws client")
		apps.Logs.Error(err)
		return nil, err
	}

	awsctx := context.TODO()
	input := &ecs.ListClustersInput{}
	return client.ListClusters(awsctx, input)
}

func ECSListContainerInstancesCmd(cfg *aws.Config, name string) (*ecs.ListContainerInstancesOutput, error) {
	if cfg == nil || cfg.Credentials == nil {
		err := errors.New("invalid aws config")
		apps.Logs.Error(err)
		return nil, err
	}

	if name == "" {
		err := errors.New("invalid arguments: empty name")
		apps.Logs.Error(err)
		return nil, err
	}

	client := ecs.NewFromConfig(*cfg)
	if client == nil {
		err := errors.New("failed to initialize aws client")
		apps.Logs.Error(err)
		return nil, err
	}

	awsctx := context.TODO()
	input := &ecs.ListContainerInstancesInput{}
	input.Cluster = &name
	return client.ListContainerInstances(awsctx, input)
}

func ECSListTaskDefinitionsCmd(cfg *aws.Config) (*ecs.ListTaskDefinitionsOutput, error) {
	if cfg == nil || cfg.Credentials == nil {
		err := errors.New("invalid aws config")
		apps.Logs.Error(err)
		return nil, err
	}

	client := ecs.NewFromConfig(*cfg)
	if client == nil {
		err := errors.New("failed to initialize aws client")
		apps.Logs.Error(err)
		return nil, err
	}

	awsctx := context.TODO()
	input := &ecs.ListTaskDefinitionsInput{}
	return client.ListTaskDefinitions(awsctx, input)
}
