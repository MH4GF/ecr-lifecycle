package ecr

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func (c *client) DescribeImages(repositoryName string) error {
	input := &ecr.DescribeImagesInput{
		RepositoryName: &repositoryName,
	}

	result, err := c.ecr.DescribeImages(input)
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}
