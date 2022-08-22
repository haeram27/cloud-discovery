package apis

import (
	"cloudisc/apps"
	"context"
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var (
	awscfg        *aws.Config
	credsFileStat fs.FileInfo
	awsCreds      CicdCreds
	lock          sync.Mutex
)

const credsFilePath = "/tmp/awsuser.json"

type CicdCreds struct {
	AccessKeyId     string `json:"AccessKeyId,omitempty"`     // mandatory
	SecretAccessKey string `json:"SecretAccessKey,omitempty"` // mandatory
	SessionToken    string `json:"SessionToken,omitempty"`
	Expiration      string `json:"Expiration,omitempty"`
	Region          string `json:"Region,omitempty"` // mandatory
	RoleArn         string `json:"RoleArn,omitempty"`
}

func AwsConfig() *aws.Config {
	return StsAssumeRoleConfigFromFile()
}

/*
/tmp/awsuser.json:

	{
		"AccessKeyId": "access_key_id",          // mandatory
		"SecretAccessKey": "secret_access_key",  // mandatory
		"SessionToken": "session_token",
		"Expiration": "expiration",
		"Region": "region",                      // mandatory
		"RoleArn": "role_arn"
	}
*/
func updateCredsFromFile() bool {
	newCredsFileStat, err := os.Stat(credsFilePath)
	if err != nil {
		apps.Logs.Error(err)
		return false
	}

	if credsFileStat == nil || credsFileStat.Size() != newCredsFileStat.Size() || credsFileStat.ModTime() != newCredsFileStat.ModTime() {
		apps.Logs.Debug("credential file is updated")
		credsFileStat = newCredsFileStat

		f, err := os.Open(credsFilePath)
		if err != nil {
			apps.Logs.Error(err)
			return false
		}
		defer f.Close()
		apps.Logs.Debug("open credential file successfully")

		j, err := ioutil.ReadAll(f)
		if err != nil {
			apps.Logs.Error(err)
			return false
		}

		awsCreds = CicdCreds{}
		err = json.Unmarshal(j, &awsCreds)
		if err != nil {
			apps.Logs.Error(err)
			return false
		}
	}

	return true
}

func DefaultConfig() (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		apps.Logs.Error(err)
		return *aws.NewConfig(), err
	}

	return cfg, nil
}

func SharedProfileConfig(profile string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		apps.Logs.Error(err)
		return *aws.NewConfig(), err
	}

	return cfg, nil
}

func StaticCredentialConfig(akid string, seckey string, token string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(akid, seckey, token)))
	if err != nil {
		apps.Logs.Error(err)
		return *aws.NewConfig(), err
	}

	return cfg, nil
}

/*
AssumeRoleXXXConfig() makes aws.Config performed assume-role process.
stsUserCfg can be made by  DefaultConfig(), SharedProfileConfig(), StaticCredentialConfig()
*/
func AssumeRoleConfig(stsUserCfg *aws.Config, roleArn string) {
	stsSvc := sts.NewFromConfig(*stsUserCfg)
	creds := stscreds.NewAssumeRoleProvider(stsSvc, roleArn)

	stsUserCfg.Credentials = aws.NewCredentialsCache(creds)
}

func AssumeRoleCustomMFAConfig(stsUserCfg *aws.Config, roleArn string, mfaSerialNumber string, mfaToken string) {
	staticTokenProvider := func() (string, error) {
		return mfaToken, nil
	}

	creds := stscreds.NewAssumeRoleProvider(sts.NewFromConfig(*stsUserCfg), roleArn, func(o *stscreds.AssumeRoleOptions) {
		o.SerialNumber = aws.String(mfaSerialNumber)
		o.TokenProvider = staticTokenProvider
	})

	stsUserCfg.Credentials = aws.NewCredentialsCache(creds)
}

/*
akId : user's Access Key ID
secKey : user's Secret Access Key
roleArn : arn of role
*/
func StsAssumeRoleConfig(c *CicdCreds) (*aws.Config, error) {
	cfg, err := StaticCredentialConfig(c.AccessKeyId, c.SecretAccessKey, "")
	if err != nil {
		apps.Logs.Error(err)
		return &aws.Config{}, err
	}

	cfg.Region = c.Region

	if c.RoleArn == "" {
		// apps.Logs.Info("no role information in credential file")
	} else {
		AssumeRoleConfig(&cfg, c.RoleArn)
	}

	return &cfg, nil
}

/*
auth for testing
*/
func StsAssumeRoleConfigFromFile() *aws.Config {
	lock.Lock()
	if updateCredsFromFile() {
		new, err := StsAssumeRoleConfig(&awsCreds)
		if err != nil {
			apps.Logs.Error(err)
			return awscfg
		}
		awscfg = new
	}
	lock.Unlock()
	return awscfg
}
