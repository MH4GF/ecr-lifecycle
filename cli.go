package main

import (
	"github.com/MH4GF/ecr-lifecycle/ecr"
	"github.com/urfave/cli/v2"
	"sync"
)

func newApp() *cli.App {
	app := cli.NewApp()

	app.Name = "ecr-lifecycle"
	app.Commands = []*cli.Command{
		&cmdDeleteImages,
	}

	return app
}

//cliのフラグをチェックし実行
var cmdDeleteImages = cli.Command{
	Name: "delete-images",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "profile",
			Aliases: []string{"p"},
			Usage:   "if development, can select specify profile from ~/.aws/credentials",
		},
		&cli.StringFlag{
			Name:     "template",
			Aliases:  []string{"t"},
			Usage:    "load YAML file for configuration.",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		// 入力されたflagを元に各種init
		config, err := newConfig(c)
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		semaphore := make(chan struct{}, 10)

		//repositoryごとに並行で回す
		for _, repo := range config.Repositories {
			wg.Add(1)
			semaphore <- struct{}{}

			go func(r ecr.Repository) {
				defer func() {
					<-semaphore
					defer wg.Done()
				}()

				result, err := config.EcrClient.BatchDeleteImages(r, config.Keep, config.EcsAllRunningTasks)
				if err != nil {
					log.sugar.Warnf("could not delete images: %s", err)
				}
				if result != nil {
					for _, f := range result.Failures {
						log.sugar.Warnw("warn", "FailureCode", f.FailureCode, "FailureReason", f.FailureReason, "ImageId", f.ImageId)
					}
					for _, id := range result.ImageIds {
						log.sugar.Infow("deletedImageId", "RepositoryName", r.Detail.RepositoryName, "ImageDigest", id.ImageDigest)
					}
				}
			}(repo)
		}
		wg.Wait()

		return nil
	},
}
