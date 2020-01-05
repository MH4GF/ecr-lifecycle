package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type Client struct {
	ecs *ecs.ECS
}

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
