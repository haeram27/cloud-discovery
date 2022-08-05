package apis

import (
	apps "awsdisc/apps"
	"awsdisc/apps/util"
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func ECRDescribeRegistryCmd(cfg *aws.Config) (*ecr.DescribeRegistryOutput, error) {
	if cfg == nil || cfg.Credentials == nil {
		err := errors.New("invalid aws config")
		apps.Logs.Error(err)
		return nil, err
	}

	client := ecr.NewFromConfig(*cfg)
	if client == nil {
		err := errors.New("failed to initialize aws client")
		apps.Logs.Error(err)
		return nil, err
	}

	awsctx := context.TODO()
	input := &ecr.DescribeRegistryInput{}
	return client.DescribeRegistry(awsctx, input)
}

func ECRDescribeRepositoriesCmd(cfg *aws.Config) (*ecr.DescribeRepositoriesOutput, error) {
	if cfg == nil || cfg.Credentials == nil {
		err := errors.New("invalid aws config")
		apps.Logs.Error(err)
		return nil, err
	}

	client := ecr.NewFromConfig(*cfg)
	if client == nil {
		err := errors.New("failed to initialize aws client")
		apps.Logs.Error(err)
		return nil, err
	}

	awsctx := context.TODO()
	input := &ecr.DescribeRepositoriesInput{}
	return client.DescribeRepositories(awsctx, input)
}

func ECRListImagesCmd(cfg *aws.Config, repoName string) (*ecr.ListImagesOutput, error) {
	if cfg == nil || cfg.Credentials == nil {
		err := errors.New("invalid aws config")
		apps.Logs.Error(err)
		return nil, err
	}

	if repoName == "" {
		err := errors.New("invalid repository name")
		apps.Logs.Error(err)
		return nil, err
	}

	client := ecr.NewFromConfig(*cfg)
	if client == nil {
		err := errors.New("failed to initialize aws client")
		apps.Logs.Error(err)
		return nil, err
	}

	awsctx := context.TODO()
	input := &ecr.ListImagesInput{
		RepositoryName: &repoName,
	}

	return client.ListImages(awsctx, input)
}

func ECRListImagesAll(cfg *aws.Config) []string {
	if cfg == nil || cfg.Credentials == nil {
		err := errors.New("invalid aws config")
		apps.Logs.Error(err)
		return nil
	}

	var jsonBlob []byte
	result, err := ECRDescribeRepositoriesCmd(cfg)
	if err != nil {
		apps.Logs.Error(err)
		return nil
	} else {
		jsonBlob, err = json.Marshal(result)
		if err != nil {
			apps.Logs.Error(err)
			return nil
		}
	}

	names := util.JsonPath(jsonBlob, "$.Repositories[:].RepositoryName")

	for _, info := range names {
		apps.Logs.Debug("============================== repository name: ", info.(string))
		ECRListImagesCmd(cfg, info.(string))
	}

	return nil
}

type EcrImage struct {
	repoUri *string
	tag     *string
	digest  *string
}

type EcrImageUri interface {
	TagUri() string
	DigestUri() string
}

func (img EcrImage) TagUri() string {
	return *img.repoUri + ":" + *img.tag
}

func (img EcrImage) DigestUri() string {
	return *img.repoUri + "@" + *img.digest
}

func ECRListImagesAllST(cfg *aws.Config) []EcrImage {
	if cfg == nil || cfg.Credentials == nil {
		err := errors.New("invalid aws config")
		apps.Logs.Error(err)
		return []EcrImage{}
	}

	reposOut, err := ECRDescribeRepositoriesCmd(cfg)
	if err != nil {
		apps.Logs.Error(err)
		return []EcrImage{}
	} else {
		jsonBlob, err := json.Marshal(reposOut)
		if err != nil {
			apps.Logs.Error(err)
			return []EcrImage{}
		}

		util.PrintPrettyJson(jsonBlob)
	}

	var images []EcrImage
	for _, repo := range reposOut.Repositories {
		imgOut, err := ECRListImagesCmd(cfg, *repo.RepositoryName)
		if err != nil {
			apps.Logs.Error(err)
			continue
		}

		for _, img := range imgOut.ImageIds {
			images = append(images, EcrImage{repo.RepositoryUri, img.ImageTag, img.ImageDigest})
		}
	}

	return images
}

/*
	ecr.GetAuthorizationToken API return token(base64 encoded string)
	that has format as <USERNAME>:<PASSWORD> when it decoded.
	<USERNAME> is "AWS" so far and <PASSWORD> is token encoded string.
	For the use of returned token, token SHOULD be base64 decoded and distinguished with <USERNAME> and <PASSWORD>.

	// EXAMPLE
	data, err := base64.URLEncoding.DecodeString(awsEcrAuthTok)
	if err != nil {
		apps.Logs.Error(err)
	}

	parts := strings.SplitN(string(data), ":", 2)

	auth := types.AuthConfig{
		Username: parts[0],
		Password: parts[1],
	}
*/
func ECRGetAuthorizationTokenCmd(cfg *aws.Config) (string, error) {
	if cfg == nil || cfg.Credentials == nil {
		err := errors.New("invalid aws config")
		apps.Logs.Error(err)
		return "", err
	}

	client := ecr.NewFromConfig(*cfg)
	if client == nil {
		err := errors.New("failed to initialize aws client")
		apps.Logs.Error(err)
		return "", err
	}

	awsctx := context.TODO()
	input := &ecr.GetAuthorizationTokenInput{}

	out, err := client.GetAuthorizationToken(awsctx, input)
	if err != nil {
		err := errors.New("failed to initialize aws client")
		apps.Logs.Error(err)
		return "", err
	}

	if out.AuthorizationData != nil && len(out.AuthorizationData) > 0 {
		return *out.AuthorizationData[0].AuthorizationToken, nil
	}

	return "", nil
}
