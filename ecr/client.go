package ecr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// Client ... Store ECR client with a session
type Client struct {
	ecr *ecr.ECR
}

// NewClient ... Create a ECR client with profile and region
func NewClient(awsProfile string, awsRegion string) *Client {
	c := &Client{}

	// インスタンスのIAMロール or ローカルのawsProfileを使ってセッションを張る
	var sess *session.Session
	if awsProfile != "" {
		sess = session.Must(session.NewSessionWithOptions(session.Options{Profile: awsProfile}))
	} else {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Config: *aws.NewConfig().WithCredentialsChainVerboseErrors(true),
		}))
	}

	config := aws.NewConfig().WithRegion(awsRegion).WithCredentialsChainVerboseErrors(true)
	c.ecr = ecr.New(sess, config)

	return c
}
