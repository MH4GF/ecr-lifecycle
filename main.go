package main

import (
	"fmt"
	"github.com/Taimee/ecr-lifecycle/ecs"
)

func run() error {
	client, err := ecs.NewClient("timee-jp-prod", "ap-northeast-1")
	if err != nil {
		return err
	}
	if _, err := client.ListAllRunningTasks(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
