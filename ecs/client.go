package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Client ... ECS client with a session
type Client struct {
	ECS *ecs.ECS
}

// NewClient is constructor
func NewClient(awsProfile string, awsRoleArn string, awsRegion string) (*Client, error) {
	c := &Client{}

	var baseSess *session.Session
	if awsProfile != "" {
		baseSess = session.Must(session.NewSessionWithOptions(session.Options{Profile: awsProfile}))
	} else {
		baseSess = session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{CredentialsChainVerboseErrors: aws.Bool(true)},
		}))
	}
	assumeRoler := sts.New(baseSess)

	creds := stscreds.NewCredentialsWithClient(assumeRoler, awsRoleArn)
	config := aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(true),
		Region:                        &awsRegion,
		Credentials:                   creds,
	}
	sess, err := session.NewSession(&config)
	if err != nil {
		return nil, err
	}

	c.ECS = ecs.New(sess, &config)

	return c, err
}
