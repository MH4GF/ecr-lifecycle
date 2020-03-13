package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// ListClusters ... クラスター一覧を取得する
func (c *Client) ListClusters() ([]*string, error) {
	var nextToken *string
	var clusterArns []*string

	for {
		input := &ecs.ListClustersInput{
			MaxResults: aws.Int64(100), // 最大値
			NextToken: nextToken,
		}
		result, err := c.ECS.ListClusters(input)
		if err != nil {
			return nil, err
		}
		clusterArns = append(clusterArns, result.ClusterArns...)

		if result.NextToken != nil {
			nextToken = result.NextToken
		} else {
			return clusterArns, nil // result.NextTokenがなければ終了
		}
	}
}
