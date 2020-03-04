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
			Name:  "ecr-profile",
			Usage: "Use a specific profile of stored ecr repository from your credential file.",
		},
		&cli.StringSliceFlag{
			Name:  "ecs-profiles",
			Usage: "Use a multiple profiles of running ecs task from your credential file.",
		},
		&cli.StringFlag{
			Name:    "region",
			Aliases: []string{"r"},
			Usage:   "The region to use.",
		},
		&cli.IntFlag{
			Name:    "keep",
			Value:   10000,
			Aliases: []string{"k"},
			Usage:   "Number of images to keep from latest.",
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
		for _, repo := range config.repositories {
			wg.Add(1)
			semaphore <- struct{}{}

			go func(r ecr.Repository) {
				defer func() {
					<-semaphore
					defer wg.Done()
				}()

				result, err := config.ecrClient.BatchDeleteImages(r, config.flag.keep, config.ecsAllRunningTasks)
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
