package ecr

import (
	"github.com/aws/aws-sdk-go/service/ecr"
)

// Repository ... Store original ecr.Repository
type Repository struct {
	Detail *ecr.Repository
}

// DescribeRepositories ... clientのprofileにひもづくECRリポジトリ一覧を取得する
func (c *Client) DescribeRepositories() ([]Repository, error) {
	var rs []Repository
	input := &ecr.DescribeRepositoriesInput{}

	result, err := c.ecr.DescribeRepositories(input)
	if err != nil {
		return rs, err
	}

	for _, repo := range result.Repositories {
		rs = append(rs, Repository{Detail: repo})
	}

	return rs, nil
}
