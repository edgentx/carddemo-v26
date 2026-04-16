.PHONY: build test lint

build:
	@echo "Building..."
	@go build ./...

test:
	@echo "Testing..."
	@go test ./... -v -cover

lint:
	@echo "Linting..."
	@golangci-lint run
