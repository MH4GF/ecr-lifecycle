package ecr

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type client struct {
	ecr *ecr.ECR
}

func NewClient(awsProfile string) (*client, error) {
	c := &client{}

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           awsProfile,
	})

	if err != nil {
		return nil, err
	}
	c.ecr = ecr.New(sess)

	return c, err
}

func (c *client) describeRepositories() error {
	input := &ecr.DescribeRepositoriesInput{}
	result, err := c.ecr.DescribeRepositories(input)
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}
