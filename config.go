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

// Config ... app実行に伴う全ての設定項目を格納
type Config struct {
	// ECSタスクが実行されているAWSアカウントにスイッチロールするためのAssumeRoleArnの配列
	EcsAssumeRoleArns []string `yaml:"EcsAssumeRoleArns"`

	// 動作するregion
	Region string `yaml:"Region"`

	// 最新から何件保持するか指定する。
	Keep int `yaml:"Keep"`

	// ECRのセッションクライアント
	EcrClient ecr.Client

	// ECRに存在する全てのリポジトリ
	Repositories []ecr.Repository

	// EcsAssumeRoleArnsで指定したAWSアカウントで実行されている全てのECSタスク
	EcsAllRunningTasks []ecs.Task
}

func newConfig(c *cli.Context) (*Config, error) {
	config := Config{}
	p := c.String("profile")
	t := c.String("template")
	if t != "" {
		if err := config.loadYaml(t); err != nil {
			return nil, errors.New(fmt.Sprintf("error on reading template file: %s", err))
		}
	} else {
		config.EcsAssumeRoleArns = c.StringSlice("ecs-assume-role-arns")
		config.Region = c.String("region")
		config.Keep = c.Int("keep")
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	// ecrClientのinit
	config.EcrClient = *ecr.NewClient(p, config.Region)

	// repositoryを取得
	repositories, err := config.EcrClient.DescribeRepositories()
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

func (c *Config) validate() error {
	// EcsAssumeRoleArns
	if len(c.EcsAssumeRoleArns) == 0 {
		return errors.New("invalid params: EcrAssumeRoleArn is required")
	}
	for _, arn := range c.EcsAssumeRoleArns {
		if arn == "" {
			return errors.New("invalid params: EcsAssumeRoleArns is required and minimum field size of 20")
		}
		if arn != "" && len(arn) < 20 {
			return errors.New("invalid params: EcsAssumeRoleArns is required and minimum field size of 20")
		}
	}

	// Region
	if c.Region == "" {
		return errors.New("invalid params: Region is required")
	}

	// Keep
	if c.Keep < 1 {
		return errors.New("invalid params: Keep is required and minimum field size of 1")
	}

	return nil
}
