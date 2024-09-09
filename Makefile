BINARY_NAME := notification-service

all: build

build:
	CGO_ENABLED=0 go build -o bin/$(BINARY_NAME) .

test:
	go test -v ./...

deps:
	go mod download

clean:
	rm -f bin/$(BINARY_NAME)

.PHONY: all build run test clean deps
