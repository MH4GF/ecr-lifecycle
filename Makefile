REVISION    := $(shell git describe --always)
LDFLAGS     := -ldflags="-X \"main.Revision=$(REVISION)\" -extldflags \"-static\""
GO_IMAGE    := golang:1.13.4

NAME        := ecr-lifecycle
linux_name	:= $(name)-linux-amd64
darwin_name	:= $(name)-darwin-amd64

.PHONY: help build clean run lint fmt test

help:
	@awk -F ':|##' '/^[^\t].+?:.*?##/ { printf "\033[36m%-22s\033[0m %s\n", $$1, $$NF }' $(MAKEFILE_LIST)

build: clean ## go build
	@go build -o bin/$(NAME) $(LDFLAGS)

build-cross: clean ## build for darwin and linux to bin/
	GOOS=linux GOARCH=amd64 go build -o bin/$(linux_name) $(LDFLAGS) *.go
	GOOS=darwin GOARCH=amd64 go build -o bin/$(darwin_name) $(LDFLAGS) *.go

docker-build: clean ## go build on docker
	@docker run --rm \
		-e "GO111MODULE=on" \
		-v `pwd`:/go/src/github.com/Taimee/$(NAME) \
		-w /go/src/github.com/Taimee/$(NAME) \
		$(GO_IMAGE) bash build.sh

clean: ## remove binary
	@rm -f bin/*

run: ## go run
	@go run main.go

lint: ## golint
	@golint $(go list ./... | grep -v /vendor/)

fmt: ## go fmt
	@go fmt ./...

test: ## go test
	@go test ./... -v
