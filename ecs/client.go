package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// Client ... ECS client with a session
type Client struct {
	ecs *ecs.ECS
}

// NewClient ... Create a ECS client with profile and region
func NewClient(awsProfile string, awsRegion string) (*Client, error) {
	c := &Client{}

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           awsProfile,
	})

	if err != nil {
		return nil, err
	}
	c.ecs = ecs.New(sess, aws.NewConfig().WithRegion(awsRegion))

	return c, err
}
