BINARY_NAME := notification-service
BINARY_NAME_CLI := notification-service-cli

all: build


build-web:
	CGO_ENABLED=0 go build -o bin/$(BINARY_NAME) cmd/web/main.go

build-cli:
	CGO_ENABLED=0 go build -o bin/$(BINARY_NAME_CLI) cmd/cli/main.go
	
build: build-web build-cli

test:
	go test -v ./...

deps:
	go mod download

clean:
	rm -f bin/{$(BINARY_NAME),$(BINARY_NAME_CLI)}

.PHONY: all build-web build-cli build run test clean deps
