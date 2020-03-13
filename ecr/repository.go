package ecr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// Repository ... Store original ecr.Repository
type Repository struct {
	Detail *ecr.Repository
}

// DescribeRepositories ... clientのprofileにひもづくECRリポジトリ一覧を取得する
func (c *Client) DescribeRepositories() ([]Repository, error) {
	var rs []Repository
	var nextToken *string

	for {
		input := &ecr.DescribeRepositoriesInput{
			MaxResults: aws.Int64(1000), // 最大値
			NextToken:  nextToken,
		}
		result, err := c.ecr.DescribeRepositories(input)
		if err != nil {
			return rs, err
		}

		for _, repo := range result.Repositories {
			rs = append(rs, Repository{Detail: repo})
		}

		if result.NextToken != nil {
			nextToken = result.NextToken
		} else {
			return rs, nil // result.NextTokenがなければ終了
		}
	}
}
