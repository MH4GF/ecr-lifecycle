#!/bin/bash -

declare -r Name="ecr-lifecycle"

for GOOS in darwin linux; do
    GO111MODULE=on GOOS=$GOOS GOARCH=amd64 go build -o bin/${Name}-${GOOS}-amd64 ./*.go
done
