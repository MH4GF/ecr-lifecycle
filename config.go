package main

import (
	"github.com/MH4GF/ecr-lifecycle/ecr"
	"github.com/MH4GF/ecr-lifecycle/ecs"
	"github.com/urfave/cli/v2"
)

// Config ... app実行に伴う全てのstructを格納
type Config struct {
	flag         Flag
	ecrClient    ecr.Client
	ecsClients   []ecs.Client
	repositories []ecr.Repository
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

	// ecsClientを全てinit
	var ecsClients []ecs.Client
	for _, p := range f.ecsProfiles {
		sess, err := ecs.RegisterECSNewSession(p, f.region)
		if err != nil {
			return nil, err
		}
		ecsClient := ecs.NewClient(sess)
		ecsClients = append(ecsClients, *ecsClient)
	}

	// repositoryを取得
	repositories, err := ecrClient.DescribeRepositories()
	if err != nil {
		return nil, err
	}
	for _, r := range repositories {
		log.sugar.Infof("target repositoryArn: %s", *r.Detail.RepositoryArn)
	}

	// 全てconfigとしてぶち込む
	config := Config{
		flag:         f,
		ecrClient:    *ecrClient,
		ecsClients:   ecsClients,
		repositories: repositories,
	}
	return &config, err
}
