package ecr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type Client struct {
	ecr *ecr.ECR
	region *string
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
	c.ecr = ecr.New(sess, aws.NewConfig().WithRegion(awsRegion))
	c.region = &awsRegion

	return c, err
}
