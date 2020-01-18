package main

import (
	"errors"
	"github.com/Taimee/ecr-lifecycle/ecr"
	"github.com/urfave/cli/v2"
	"strconv"
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
			Usage:   "Use a specific profile from your credential file.",
		},
		&cli.StringFlag{
			Name:    "region",
			Aliases: []string{"r"},
			Usage:   "The region to use.",
		},
		&cli.StringFlag{
			Name:    "keep",
			Value:   "10000",
			Aliases: []string{"k"},
			Usage:   "Number of images to keep from latest.",
		},
	},
	Action: func(c *cli.Context) error {
		// Get flag options
		profile := c.String("profile")
		region := c.String("region")
		keep := c.String("keep")

		// Flag check
		if profile == "" {
			return errors.New("-p, --profile option is required")
		}
		if region == "" {
			return errors.New("-r or --region option is required")
		}

		client, err := ecr.NewClient(profile, region)
		if err != nil {
			return err
		}

		repositories, err := client.DescribeRepositories()
		if err != nil {
			return err
		}
		for _, r := range repositories {
			log.sugar.Infof("target repositoryArn: %s", *r.Detail.RepositoryArn)
		}

		num, err := strconv.Atoi(keep)
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		semaphore := make(chan struct{}, 10)

		for _, repo := range repositories {
			wg.Add(1)
			semaphore <- struct{}{}

			go func(r ecr.Repository) {
				defer func() {
					<-semaphore
					defer wg.Done()
				}()
				result, err := client.BatchDeleteImages(r, &num)
				if err != nil {
					log.sugar.Warnf("could not delete images: %s", err)
				}
				if result != nil {
					for _, f := range result.Failures {
						log.sugar.Warnw("warn", "FailureCode", f.FailureCode, "FailureReason", f.FailureReason, "ImageId", f.ImageId)
					}
					for _, id := range result.ImageIds {
						log.sugar.Infow("deletedImageId", "ImageDigest", id.ImageDigest)
					}
				}
			}(repo)
		}
		wg.Wait()

		return nil
	},
}
