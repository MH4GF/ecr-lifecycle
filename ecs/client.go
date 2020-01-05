package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// ECS ... モックに差し替えやすくするため, aws-sdkのメソッドをinterfaceにしている
type ECS interface {
	DescribeTaskDefinition(input *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error)
	DescribeTasks(input *ecs.DescribeTasksInput) (*ecs.DescribeTasksOutput, error)
	ListTasks(input *ecs.ListTasksInput) (*ecs.ListTasksOutput, error)
	ListClusters(input *ecs.ListClustersInput) (*ecs.ListClustersOutput, error)
}

// Client ... ECS client with a session
type Client struct {
	ecs ECS
}

// RegisterECSNewSession ... Create a ECS client with profile and region
func RegisterECSNewSession(awsProfile string, awsRegion string) (*ecs.ECS, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           awsProfile,
	})

	if err != nil {
		return nil, err
	}
	ecsClient := ecs.New(sess, aws.NewConfig().WithRegion(awsRegion))

	return ecsClient, nil
}

// NewClient is constructor
func NewClient(ecs ECS) *Client {
	return &Client{ecs}
}
