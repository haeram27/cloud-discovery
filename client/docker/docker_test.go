package dkr

import (
	aws "awsdisc/client/aws/apis"
	"testing"
)

func TestLogin(t *testing.T) {
	tok, err := aws.ECRGetAuthorizationTokenCmd(aws.AwsConfig())
	if err != nil {
		t.Fatal(err)
	}

	Login("797216966998.dkr.ecr.ap-northeast-2.amazonaws.com", tok)
}

func TestPullImage(t *testing.T) {
	tok, err := aws.ECRGetAuthorizationTokenCmd(aws.AwsConfig())
	if err != nil {
		t.Fatal(err)
	}

	PullImage("797216966998.dkr.ecr.ap-northeast-2.amazonaws.com/cicd-sec/alpine@sha256:4ff3ca91275773af45cb4b0834e12b7eb47d1c18f770a0b151381cd227f4c253", tok)
}
