package ecr

import (
	"github.com/aws/aws-sdk-go/service/ecr"
)

type repository struct {
	original *ecr.Repository
}

func (c *Client) DescribeRepositories() ([]*repository, error) {
	input := &ecr.DescribeRepositoriesInput{}

	result, err := c.ecr.DescribeRepositories(input)
	if err != nil {
		return nil, err
	}

	var repositories []*repository
	for _, repo := range result.Repositories {
		repositories = append(repositories, &repository{original: repo})
	}

	return repositories, nil
}
