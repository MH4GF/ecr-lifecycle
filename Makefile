NAME     := ecr-lifecycle
REVISION := $(shell git describe --always)
LDFLAGS  := -ldflags="-X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

.PHONY: help build build-image clean run lint fmt test


help:
	@awk -F ':|##' '/^[^\t].+?:.*?##/ { printf "\033[36m%-22s\033[0m %s\n", $$1, $$NF }' $(MAKEFILE_LIST)

build:
	@go build -o bin/$(NAME) $(LDFLAGS)

build-image:
	@docker build . -t $(NAME)
clean: ## remove binary
	@rm -f bin/$(NAME)


run:
	@go run main.go

lint:
	@golint $(go list ./... | grep -v /vendor/)

fmt:
	@go fmt ./...

test:
	@go test ./... -v
