package ecr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Client ... Store ECR client with a session
type Client struct {
	ecr *ecr.ECR
}

// NewClient ... Create a ECR client with profile and region
func NewClient(awsProfile string, awsRoleArn string, awsRegion string) (*Client, error) {
	c := &Client{}

	// インスタンスのIAMロール or ローカルのawsProfileを使ってまずセッションを張る
	var baseSess *session.Session
	if awsProfile != "" {
		baseSess = session.Must(session.NewSessionWithOptions(session.Options{Profile: awsProfile}))
	} else {
		baseSess = session.Must(session.NewSessionWithOptions(session.Options{
			Config: *aws.NewConfig().WithCredentialsChainVerboseErrors(true),
		}))
	}
	assumeRoler := sts.New(baseSess)

	// 指定したECRへassumeRole
	creds := stscreds.NewCredentialsWithClient(assumeRoler, awsRoleArn)
	config := aws.NewConfig().WithRegion(awsRegion).WithCredentials(creds).WithCredentialsChainVerboseErrors(true)
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	c.ecr = ecr.New(sess, config)

	return c, err
}
