package apis

import (
	"cloudisc/apps"
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic"
)

func ECRPubDescribeRegistryCmd(cfg *aws.Config) (*ecrpublic.DescribeRegistriesOutput, error) {
	if cfg == nil || cfg.Credentials == nil {
		err := errors.New("invalid aws config")
		apps.Logs.Error(err)
		return nil, err
	}

	client := ecrpublic.NewFromConfig(*cfg)
	if client == nil {
		err := errors.New("failed to initialize aws client")
		apps.Logs.Error(err)
		return nil, err
	}

	awsctx := context.TODO()
	input := &ecrpublic.DescribeRegistriesInput{}
	return client.DescribeRegistries(awsctx, input)
}
