# SB Media Processor Makefile

# Metadata
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build variables
BINARY_NAME = sb
MAIN_PKG = .
BUILD_DIR = dist
LDFLAGS = -ldflags="-s -w -X 'github.com/onedusk/sb/cmd.Version=$(VERSION)' -X 'github.com/onedusk/sb/cmd.Commit=$(COMMIT)' -X 'github.com/onedusk/sb/cmd.Date=$(DATE)'"

# Go commands
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOMOD = $(GOCMD) mod

# Detect OS
UNAME_S := $(shell uname -s)

.PHONY: all build clean test deps install darwin linux windows help

# Default target
all: clean deps test build

# Build for current platform
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PKG)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for macOS (darwin)
darwin:
	@echo "Building for macOS (darwin/amd64 and darwin/arm64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PKG)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PKG)
	@echo "Darwin builds complete"

# Build for Linux
linux:
	@echo "Building for Linux (linux/amd64 and linux/arm64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PKG)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PKG)
	@echo "Linux builds complete"

# Build for Windows
windows:
	@echo "Building for Windows (windows/amd64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PKG)
	@echo "Windows build complete"

# Build for all platforms
release: clean deps test darwin linux windows
	@echo "All platform builds complete"
	@ls -lh $(BUILD_DIR)/

# Install binary to system
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
ifeq ($(UNAME_S),Darwin)
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
else ifeq ($(UNAME_S),Linux)
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
endif
	@echo "Installation complete: /usr/local/bin/$(BINARY_NAME)"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies updated"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "Tests complete"

# Run tests with coverage report
test-coverage: test
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	$(GOCLEAN)
	@echo "Clean complete"

# Run the binary
run: build
	$(BUILD_DIR)/$(BINARY_NAME)

# Display help
help:
	@echo "SB Media Processor Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all           - Clean, download deps, test, and build"
	@echo "  build         - Build for current platform"
	@echo "  darwin        - Build for macOS (amd64 + arm64)"
	@echo "  linux         - Build for Linux (amd64 + arm64)"
	@echo "  windows       - Build for Windows (amd64)"
	@echo "  release       - Build for all platforms"
	@echo "  install       - Install binary to /usr/local/bin"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  test          - Run tests with race detector"
	@echo "  test-coverage - Run tests and generate coverage report"
	@echo "  clean         - Remove build artifacts"
	@echo "  run           - Build and run the binary"
	@echo "  help          - Display this help message"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION       - Version string (default: git describe or 'dev')"
	@echo "  COMMIT        - Git commit hash (default: git rev-parse)"
	@echo "  DATE          - Build date (default: current UTC time)"
