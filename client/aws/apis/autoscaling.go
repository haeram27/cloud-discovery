package apis

import (
	"cloudisc/apps"
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	ascl "github.com/aws/aws-sdk-go-v2/service/autoscaling"
)

func ASCLDescribeAutoScalingGroupsCmd(cfg *aws.Config) (*ascl.DescribeAutoScalingGroupsOutput, error) {
	if cfg == nil || cfg.Credentials == nil {
		err := errors.New("invalid aws config")
		apps.Logs.Error(err)
		return nil, err
	}

	client := ascl.NewFromConfig(*cfg)
	if client == nil {
		err := errors.New("failed to initialize aws client")
		apps.Logs.Error(err)
		return nil, err
	}

	awsctx := context.TODO()
	input := &ascl.DescribeAutoScalingGroupsInput{}
	return client.DescribeAutoScalingGroups(awsctx, input)
}
