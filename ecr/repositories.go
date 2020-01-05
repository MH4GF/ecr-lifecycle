package ecr

import (
	"github.com/aws/aws-sdk-go/service/ecr"
)

// Repository is stored original ecr.Repository
type Repository struct {
	original *ecr.Repository
}

// DescribeRepositories ... clientのprofileにひもづくECRリポジトリ一覧を取得する
func (c *Client) DescribeRepositories() ([]*Repository, error) {
	input := &ecr.DescribeRepositoriesInput{}

	result, err := c.ecr.DescribeRepositories(input)
	if err != nil {
		return nil, err
	}

	var repositories []*Repository
	for _, repo := range result.Repositories {
		repositories = append(repositories, &Repository{original: repo})
	}

	return repositories, nil
}
