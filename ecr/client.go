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

	return c, err
}
