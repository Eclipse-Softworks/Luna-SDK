.PHONY: all build test lint fmt clean

all: lint test build

build:
	go build -v ./...

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

coverage: test
	go tool cover -html=coverage.txt -o coverage.html

lint:
	golangci-lint run

fmt:
	gofmt -s -w .
	goimports -w -local github.com/eclipse-softworks/luna-sdk-go .

clean:
	go clean
	rm -f coverage.txt coverage.html

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
