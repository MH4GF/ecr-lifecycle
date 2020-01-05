package ecs

import "github.com/aws/aws-sdk-go/service/ecs"

// ListClusters ... クラスター一覧を取得する
func (c *Client) ListClusters() (*ecs.ListClustersOutput, error) {
	result, err := c.ecs.ListClusters(&ecs.ListClustersInput{})
	if err != nil {
		return nil, err
	}

	return result, nil
}
