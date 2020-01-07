.PHONY: help build build-image clean run lint fmt test

name := ecr-lifecycle

help:
	@awk -F ':|##' '/^[^\t].+?:.*?##/ { printf "\033[36m%-22s\033[0m %s\n", $$1, $$NF }' $(MAKEFILE_LIST)

build:
	@go build -o $(name)

build-image:
	@docker build . -t $(name)

clean:
	@rm -f $(name)

run:
	@go run main.go

lint:
	@golint $(go list ./... | grep -v /vendor/)

fmt:
	@go fmt ./...

test:
	@go test ./... -v
