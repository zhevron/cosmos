.DEFAULT_GOAL := build

.PHONY:build install-tools lint test

build:
	go build ./...

install-tools:
	GO111MODULE=off go get \
		github.com/jstemmer/go-junit-report

lint:
	golangci-lint run ./...

test:
	go test -v ./...
	go test -v ./... 2>&1 | go-junit-report > report.xml
