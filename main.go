package main

import (
	"fmt"
	"github.com/Taimee/ecr-lifecycle/ecr"
	"os"
)

func run() error {
	client, err := ecr.NewClient("sandbox", "ap-northeast-1")
	if err != nil {
		return err
	}

	repositories, err := client.DescribeRepositories()
	if err != nil {
		return err
	}

	count := 0
	for _, repo := range repositories {
		client.BatchDeleteImages(repo, &count)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
