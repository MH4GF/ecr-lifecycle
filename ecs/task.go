package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// Task ... 現在実行中のタスクの情報を持つ。imageがあれば良い
type Task struct {
	TaskArn string
	Image   string
}

// ListAllRunningTasks ... 現在実行中のタスク一覧を取得する
func (c *Client) ListAllRunningTasks() ([]Task, error) {
	var tasks []Task

	// 現在実行中のタスク一覧を取得
	outputs, err := c.listAllTasksOutput()
	if err != nil {
		return nil, err
	}

	for _, o := range outputs {

		// TaskArnsが空配列の場合スキップ
		if len(o.listTasksOutput.TaskArns) == 0 {
			continue
		}

		// TaskArnsを元にタスクの詳細を取得
		tasksOutput, err := c.describeTasks(o.clusterArn, o.listTasksOutput.TaskArns)
		if err != nil {
			return nil, err
		}

		// タスク詳細を元にtaskDefinitionを取得
		for _, taskOutput := range tasksOutput.Tasks {
			taskDefinition, err := c.describeTaskDefinition(taskOutput.TaskDefinitionArn)
			if err != nil {
				return nil, err
			}

			// 最終的にcontainerDefinitionからimageArnを取得する
			for _, c := range taskDefinition.TaskDefinition.ContainerDefinitions {
				tasks = append(tasks, Task{
					TaskArn: aws.StringValue(taskOutput.TaskArn),
					Image:   aws.StringValue(c.Image),
				})
			}

		}
	}

	return tasks, nil
}

func (c *Client) describeTaskDefinition(taskDefinition *string) (*ecs.DescribeTaskDefinitionOutput, error) {
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: taskDefinition,
	}

	result, err := c.ECS.DescribeTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// taskの詳細情報を取得
func (c *Client) describeTasks(clusterArn *string, taskArns []*string) (*ecs.DescribeTasksOutput, error) {
	input := &ecs.DescribeTasksInput{
		Cluster: clusterArn,
		Tasks:   taskArns,
	}

	result, err := c.ECS.DescribeTasks(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// clusterArnも持たせるためにtypeを拡張
type listTasksOutput struct {
	listTasksOutput *ecs.ListTasksOutput
	clusterArn      *string
}

// statusがRUNNINGのecs task一覧を取得
func (c *Client) listAllTasksOutput() ([]*listTasksOutput, error) {
	var tasks []*listTasksOutput

	clusterArns, err := c.ListClusters()
	if err != nil {
		return nil, err
	}

	var nextToken *string
	for _, clusterArn := range clusterArns {
		for {
			input := &ecs.ListTasksInput{
				Cluster:       clusterArn,
				DesiredStatus: aws.String("RUNNING"),
				MaxResults:    aws.Int64(100), // 最大値
				NextToken:     nextToken,
			}
			result, err := c.ECS.ListTasks(input)
			if err != nil {
				return nil, err
			}

			tasks = append(tasks, &listTasksOutput{listTasksOutput: result, clusterArn: clusterArn})
			if result.NextToken != nil {
				nextToken = result.NextToken
			} else {
				break // result.NextTokenがなければ終了
			}
		}
	}

	return tasks, nil
}
