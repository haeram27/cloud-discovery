package dkr

import (
	"awsdisc/apps"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

/*
   WARNING: This Login function to Registry is just EXAMPLE to check usage.
   Docker client can not pull image from registry without auth config in ImagePull() API
   even after this Login is successed
*/
func Login(url string, awsEcrAuthTok string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		apps.Logs.Error(err)
	}

	data, err := base64.URLEncoding.DecodeString(awsEcrAuthTok)
	if err != nil {
		apps.Logs.Error(err)
	}

	parts := strings.SplitN(string(data), ":", 2)

	auth := types.AuthConfig{
		Username:      parts[0],
		Password:      parts[1],
		ServerAddress: url,
	}

	body, err := cli.RegistryLogin(ctx, auth)
	if err != nil {
		apps.Logs.Error(err)
	}
	fmt.Printf("================== %+v", body)
}

/*
   WARNING: authentication information SHOULD be set Options of API
*/
func PullImage(uri string, awsEcrAuthTok string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		apps.Logs.Error(err)
		return err
	}

	data, err := base64.URLEncoding.DecodeString(awsEcrAuthTok)
	if err != nil {
		apps.Logs.Error(err)
		return err
	}

	parts := strings.SplitN(string(data), ":", 2)

	auth := types.AuthConfig{
		Username: parts[0],
		Password: parts[1],
	}
	authBytes, _ := json.Marshal(auth)
	authBase64 := base64.URLEncoding.EncodeToString(authBytes)

	reader, err := cli.ImagePull(ctx, uri, types.ImagePullOptions{RegistryAuth: authBase64})
	if err != nil {
		apps.Logs.Error(err)
		return err
	}

	defer reader.Close()
	io.Copy(os.Stdout, reader)

	return nil
}
