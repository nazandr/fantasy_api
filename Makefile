.PHONY: build
build:
	go build -v ./cmd/api_server

.DEFAULT_GOAL := build