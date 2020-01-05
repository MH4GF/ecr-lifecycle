package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"testing"
)

type mockedECS struct{}

func (m mockedECS) DescribeTaskDefinition(input *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
	panic("implement me")
}

func (m mockedECS) DescribeTasks(input *ecs.DescribeTasksInput) (*ecs.DescribeTasksOutput, error) {
	panic("implement me")
}

func (m mockedECS) ListTasks(input *ecs.ListTasksInput) (*ecs.ListTasksOutput, error) {
	panic("implement me")
}

func (m mockedECS) ListClusters(input *ecs.ListClustersInput) (*ecs.ListClustersOutput, error) {
	return &ecs.ListClustersOutput{
		ClusterArns: []*string{
			aws.String("arn:aws:ecs:ap-northeast-1:1234567890:cluster/hoge"),
			aws.String("arn:aws:ecs:ap-northeast-1:1234567890:cluster/fuga"),
		},
	}, nil
}

func TestListClusters(t *testing.T) {
	mockClient := NewClient(mockedECS{})
	_, err := mockClient.ListClusters()

	if err != nil {
		t.Errorf("Expected no error, but got %v.", err)
	}
}
