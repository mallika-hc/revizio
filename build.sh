#!/bin/sh

GOOS=linux GOARCH=amd64 go build
mv revizio revizio_linux_amd64

GOOS=darwin GOARCH=arm64 go build
mv revizio revizio_darwin_arm64

