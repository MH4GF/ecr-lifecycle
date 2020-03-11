package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// Client ... ECS client with a session
type Client struct {
	ECS *ecs.ECS
}

// NewClient is constructor
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
	c.ECS = ecs.New(sSess, &config)

	return c, err
}
