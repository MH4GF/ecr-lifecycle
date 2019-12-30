package main

import (
	"fmt"
	"github.com/Taimee/ecr-lifecycle/ecr"
)

func run() error {
	client, err := ecr.NewClient("timee-jp-prod")
	if err != nil {
		return err
	}

	if err = client.describeRepositories(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
