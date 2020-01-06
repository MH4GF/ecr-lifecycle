package ecr

import (
	"github.com/aws/aws-sdk-go/service/ecr"
)

// DescribeRepositories ... clientのprofileにひもづくECRリポジトリ一覧を取得する
func (c *Client) DescribeRepositories() ([]ecr.Repository, error) {
	var rs []ecr.Repository
	input := &ecr.DescribeRepositoriesInput{}

	result, err := c.ecr.DescribeRepositories(input)
	if err != nil {
		return rs, err
	}

	for _, repo := range result.Repositories {
		rs = append(rs, *repo)
	}

	return rs, nil
}
