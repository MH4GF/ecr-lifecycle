package main

func main() {
	client, err := newClient("timee-jp-prod")
	if err != nil {
		panic(err)
	}

	client.describeRepositories()
}
