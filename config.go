package main

import (
	"github.com/MH4GF/ecr-lifecycle/ecr"
	"github.com/MH4GF/ecr-lifecycle/ecs"
	"github.com/urfave/cli/v2"
)

// Config ... app実行に伴う全てのstructを格納
type Config struct {
	flag               Flag
	ecrClient          ecr.Client
	repositories       []ecr.Repository
	ecsAllRunningTasks []ecs.Task
}

func newConfig(c *cli.Context) (*Config, error) {
	// Flagの確認とcheck
	f := Flag{
		ecrProfile:  c.String("ecr-profile"),
		ecsProfiles: c.StringSlice("ecs-profiles"),
		region:      c.String("region"),
		keep:        c.Int("keep"),
	}
	if err := f.validate(); err != nil {
		return nil, err
	}

	// ecrClientのinit
	ecrClient, err := ecr.NewClient(f.ecrProfile, f.region)
	if err != nil {
		return nil, err
	}

	// repositoryを取得
	repositories, err := ecrClient.DescribeRepositories()
	if err != nil {
		return nil, err
	}
	for _, r := range repositories {
		log.sugar.Infof("target repositoryArn: %s", *r.Detail.RepositoryArn)
	}

	// ecsで現在実行しているタスクを取得
	var ecsAllRunningTasks []ecs.Task
	for _, p := range f.ecsProfiles {
		sess, err := ecs.RegisterECSNewSession(p, f.region)
		if err != nil {
			return nil, err
		}
		ecsClient := ecs.NewClient(*sess)

		tasks, err := ecsClient.ListAllRunningTasks()
		if err != nil {
			return nil, err
		}
		for _, task := range tasks {
			log.sugar.Infow("running task", "ecsAwsProfile", p, "ecsTaskArn", task.TaskArn, "taskImageUri", task.Image)
		}

		ecsAllRunningTasks = append(ecsAllRunningTasks, tasks...)
	}

	// 全てconfigとしてぶち込む
	config := Config{
		flag:               f,
		ecrClient:          *ecrClient,
		repositories:       repositories,
		ecsAllRunningTasks: ecsAllRunningTasks,
	}
	return &config, err
}
