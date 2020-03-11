package main

import (
	"errors"
	"fmt"
	"github.com/MH4GF/ecr-lifecycle/ecr"
	"github.com/MH4GF/ecr-lifecycle/ecs"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config ... app実行に伴う全てのstructを格納
type Config struct {
	EcrAssumeRoleArn   string   `yaml:"EcrAssumeRoleArn"`
	EcsAssumeRoleArns  []string `yaml:"EcsAssumeRoleArns"`
	Region             string   `yaml:"Region"`
	Keep               int      `yaml:"Keep"`
	EcrClient          ecr.Client
	Repositories       []ecr.Repository
	EcsAllRunningTasks []ecs.Task
}

func newConfig(c *cli.Context) (*Config, error) {
	p := c.String("profile")
	t := c.String("template")
	config := Config{}
	if err := config.loadYaml(t); err != nil {
		return nil, errors.New(fmt.Sprintf("error on reading template file: %s", err))
	}

	// ecrClientのinit
	ecrClient, err := ecr.NewClient(p, config.EcrAssumeRoleArn, config.Region)
	if err != nil {
		return nil, err
	}
	config.EcrClient = *ecrClient

	// repositoryを取得
	repositories, err := ecrClient.DescribeRepositories()
	if err != nil {
		return nil, err
	}
	for _, r := range repositories {
		log.sugar.Infof("target repositoryArn: %s", *r.Detail.RepositoryArn)
	}
	config.Repositories = repositories

	// ecsで現在実行しているタスクを取得
	var ecsAllRunningTasks []ecs.Task
	for _, arn := range config.EcsAssumeRoleArns {
		ecsClient, err := ecs.NewClient(p, arn, config.Region)
		if err != nil {
			return nil, err
		}

		tasks, err := ecsClient.ListAllRunningTasks()
		if err != nil {
			return nil, err
		}
		for _, task := range tasks {
			log.sugar.Infow("running task", "ecsTaskArn", task.TaskArn, "taskImageUri", task.Image)
		}

		ecsAllRunningTasks = append(ecsAllRunningTasks, tasks...)
	}
	config.EcsAllRunningTasks = ecsAllRunningTasks

	return &config, err
}

func (c *Config) loadYaml(filepath string) error {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Yaml ファイル → Struct へのパース
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return err
	}

	return nil
}
