# Development Guide

This guide covers local development setup, workflows, and best practices for contributing to SB.

## Prerequisites

- **Go 1.21+**
- **ffmpeg** (for testing actual conversions)
- **Make** (for build automation)
- **golangci-lint** (for linting)
- **Git** (for version control)

## Initial Setup

```bash
# Clone repository
git clone https://github.com/dusk-labs/sb.git
cd sb

# Install dependencies
go mod download

# Build
make build

# Run tests
make test

# Install locally
make install
```

## Project Structure

```
sb/
├── cmd/                  # CLI commands (Cobra)
│   ├── root.go          # Root command
│   ├── formats/         # Format commands
│   └── *.go             # Utility commands
├── internal/            # Internal packages
│   ├── converter/       # Converter interface & registry
│   ├── processors/      # Converter implementations
│   ├── executor/        # FFmpeg wrapper & worker pool
│   ├── config/          # Viper configuration
│   └── ui/              # Progress & output
├── docs/                # Documentation
├── .github/             # GitHub templates & workflows
├── Makefile             # Build automation
├── go.mod/go.sum        # Go dependencies
└── README.md            # User documentation
```

## Development Workflow

### 1. Create Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes

Follow coding standards:
- Run `gofmt -s -w .` before committing
- Add tests for new functionality
- Update documentation

### 3. Run Checks

```bash
# Run tests
make test

# Run linter
golangci-lint run ./...

# Check formatting
gofmt -l .

# Tidy dependencies
go mod tidy
```

### 4. Commit Changes

Use conventional commits:
```bash
git commit -m "feat(converter): add HEIC to JPG converter"
git commit -m "fix(mp4): handle special characters in filenames"
git commit -m "docs(readme): update installation instructions"
```

### 5. Push and Create PR

```bash
git push origin feature/your-feature-name
# Create PR on GitHub
```

## Testing

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/converter

# With coverage
go test -cover ./...

# With race detector
go test -race ./...

# Verbose output
go test -v ./...
```

### Writing Tests

```go
func TestConverterValidate(t *testing.T) {
    conv := NewMP4Converter()

    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid file", "test.mov", false},
        {"invalid ext", "test.txt", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := conv.Validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Debugging

### Verbose Mode

```bash
# Enable verbose logging
sb mp4 -v input.mov

# Dry run to see what would execute
sb mp4 -n -v input.mov
```

### Using Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug . -- mp4 input.mov
```

### Logging

Add debug logging in code:

```go
if opts.Verbose {
    fmt.Printf("[DEBUG] Converting: %s -> %s\n", input, output)
}
```

## Build Targets

```bash
# Build for current platform
make build

# Build for all platforms
make release

# Build for specific platform
GOOS=linux GOARCH=amd64 make build

# Install to /usr/local/bin
make install

# Clean build artifacts
make clean
```

## Code Style

### Go Standards

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Follow [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` or `gofumpt`
- Keep functions focused and small

### Project Conventions

- Package names: lowercase, single word
- Interface names: Noun (Converter, Executor)
- Constructor functions: New{Type}()
- Error messages: lowercase, no period
- Comments: Start with name of thing being described

### Example

```go
// FFmpeg wraps ffmpeg command execution
type FFmpeg struct {
    binaryPath string
}

// NewFFmpeg creates a new FFmpeg executor
func NewFFmpeg() (*FFmpeg, error) {
    path, err := exec.LookPath("ffmpeg")
    if err != nil {
        return nil, fmt.Errorf("ffmpeg not found: %w", err)
    }
    return &FFmpeg{binaryPath: path}, nil
}
```

## Adding New Features

### New Converter

See [converters.md](converters.md) for detailed guide.

Quick checklist:
1. Create processor package
2. Implement Converter interface
3. Add options struct
4. Register in init()
5. Create CLI command
6. Add tests
7. Update documentation

### New CLI Command

1. Create file in `cmd/`
2. Define Cobra command
3. Add to root in init()
4. Implement RunE function
5. Add tests

### New Configuration Option

1. Update `config.Config` struct
2. Set default in `setDefaults()`
3. Add flag in CLI command
4. Bind flag to Viper
5. Document in README

## Performance Profiling

### CPU Profiling

```bash
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```

### Memory Profiling

```bash
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

### Benchmarks

```go
func BenchmarkConvert(b *testing.B) {
    conv := NewMP4Converter()
    for i := 0; i < b.N; i++ {
        conv.Convert("test.mov", opts)
    }
}
```

## Release Process

1. Update CHANGELOG.md
2. Update version in code
3. Create git tag: `git tag -a v0.2.0 -m "Release v0.2.0"`
4. Push tag: `git push origin v0.2.0`
5. GitHub Actions builds and creates release
6. Verify release artifacts

## Troubleshooting

### Tests Failing

```bash
# Clean and rebuild
make clean && make test

# Check for stale files
go clean -cache

# Verbose test output
go test -v ./...
```

### Build Issues

```bash
# Update dependencies
go mod download
go mod tidy

# Check Go version
go version  # Should be 1.21+
```

### ffmpeg Not Found in Tests

```bash
# Install ffmpeg
brew install ffmpeg  # macOS
apt install ffmpeg   # Linux

# Verify installation
which ffmpeg
ffmpeg -version
```

## Useful Commands

```bash
# Find TODO comments
grep -r "TODO" --include="*.go" .

# Count lines of code
find . -name "*.go" | xargs wc -l

# Check for unused dependencies
go mod tidy

# Update all dependencies
go get -u ./...

# Generate documentation
godoc -http=:6060
```

## IDE Setup

### VS Code

Recommended extensions:
- Go (golang.go)
- Go Test Explorer
- golangci-lint

### GoLand

Works out of the box. Enable:
- gofmt on save
- golangci-lint inspection
- Code coverage highlighting

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Cobra Guide](https://github.com/spf13/cobra/blob/master/user_guide.md)
- [Viper Guide](https://github.com/spf13/viper#putting-values-into-viper)
- [FFmpeg Documentation](https://ffmpeg.org/documentation.html)
- [Project Issues](https://github.com/dusk-labs/sb/issues)

## Getting Help

- Check [existing issues](https://github.com/dusk-labs/sb/issues)
- Read [architecture docs](architecture.md)
- Ask in [GitHub Discussions](https://github.com/dusk-labs/sb/discussions)
- Review [ADR directory](adr/) for design decisions
