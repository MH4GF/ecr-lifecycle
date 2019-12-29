package main

import "fmt"

func run() error {
	client, err := newClient("timee-jp-prod")
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
