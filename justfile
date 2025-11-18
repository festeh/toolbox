# justfile for toolbox

# Default recipe - show available commands
default:
    @just --list

# Build the application
build:
    go build -o toolbox .

# Run the application
run:
    go run .

# Build and run
br: build
    ./toolbox

# Run tests
test:
    go test ./...

# Run tests with verbose output
test-v:
    go test -v ./...

# Format code
fmt:
    go fmt ./...

# Run go vet
vet:
    go vet ./...

# Tidy go modules
tidy:
    go mod tidy

# Download dependencies
deps:
    go mod download

# Clean build artifacts
clean:
    rm -f toolbox

# Install the binary to $GOPATH/bin
install:
    go install .

# Run all checks (fmt, vet, test)
check: fmt vet test

# Show module dependencies
mod-graph:
    go mod graph

# Update all dependencies
update:
    go get -u ./...
    go mod tidy

# Build with race detector
build-race:
    go build -race -o toolbox .

# Run with race detector
run-race:
    go run -race .
