package ecr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// Client ... Store ECR client with a session
type Client struct {
	ecr *ecr.ECR
}

// NewClient ... Create a ECR client with profile and region
func NewClient(awsProfile string, awsRoleArn string, awsRegion string) (*Client, error) {
	c := &Client{}

	baseSess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           awsProfile,
	})
	if err != nil {
		return nil, err
	}

	creds := stscreds.NewCredentials(baseSess, awsRoleArn)
	config := aws.Config{Region: &awsRegion, Credentials: creds}
	sSess, err := session.NewSession(&config)
	if err != nil {
		return nil, err
	}
	c.ecr = ecr.New(sSess, &config)

	return c, err
}
