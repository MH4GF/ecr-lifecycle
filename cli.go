package main

import (
	"errors"
	"fmt"
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
		fmt.Println(repositories)

		num, err := strconv.Atoi(keep)
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		for _, repo := range repositories {
			wg.Add(1)
			go func(r ecr.Repository) {
				defer wg.Done()
				client.BatchDeleteImages(r, &num)
			}(repo)
		}
		wg.Wait()

		return nil
	},
}
