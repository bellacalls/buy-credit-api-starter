.PHONY: run build test clean deps

# Run the application
run:
	go run cmd/api/main.go

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -rf bin/

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air
