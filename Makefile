.PHONY: build install fmt install lint test test-unit install-ci clean watch test-race test-integration release
VERSION := $(shell cat VERSION)
BIN=$(shell pwd)/bin

all: install verify

install:
	mkdir -p bin
	GO111MODULE=on GOBIN=$(BIN) go get github.com/githubnemo/CompileDaemon@v.1.1.0
	GO111MODULE=on GOBIN=$(BIN) go get github.com/giantswarm/semver-bump
	GO111MODULE=on GOBIN=$(BIN) go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.17.1
	GO111MODULE=on go mod download
	GO111MODULE=on go mod vendor
	GO111MODULE=on go mod tidy

test:
	go test ./... -timeout 120s -count 1

fmt:
	gofmt -w=true -s $$(find . -type f -name '*.go' -not -path "./vendor/*")
	goimports -w=true -d $$(find . -type f -name '*.go' -not -path "./vendor/*")

lint-code:
	./bin/golangci-lint run ./... --skip-dirs vendor

watch:
	./bin/CompileDaemon -color=true -exclude-dir=.git -build="make test"

verify: fmt lint-code test