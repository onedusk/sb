# Contributing to SB

Thank you for your interest in contributing to SB! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Adding New Converters](#adding-new-converters)

## Code of Conduct

This project adheres to a Code of Conduct. By participating, you are expected to uphold this code. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before contributing.

## Getting Started

### Prerequisites

- **Go 1.21+** installed
- **ffmpeg** installed and in PATH
- **Git** for version control
- **Make** for build automation

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/sb.git
   cd sb
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/dusk-labs/sb.git
   ```

## Development Setup

```bash
# Install dependencies
go mod download

# Build the project
make build

# Run tests
make test

# Run linter
golangci-lint run ./...

# Install the binary locally
make install
```

### Recommended Tools

- **golangci-lint**: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- **gopls**: Go language server for IDE support
- **gofumpt**: Stricter formatting than gofmt

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/dusk-labs/sb/issues)
2. If not, create a new issue using the **Bug Report** template
3. Include:
   - Clear description of the problem
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, Go version, ffmpeg version)
   - Relevant logs (use `-v` flag for verbose output)

### Suggesting Features

1. Check existing [Issues](https://github.com/dusk-labs/sb/issues) and [Discussions](https://github.com/dusk-labs/sb/discussions)
2. Create a new issue using the **Feature Request** template
3. Describe:
   - Use case and motivation
   - Proposed solution
   - Alternatives considered
   - Implementation complexity (if known)

### Requesting New Converters

Use the **Converter Request** issue template to suggest new format converters.

## Coding Standards

### Go Style Guide

Follow standard Go conventions:

- Use `gofmt` or `gofumpt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Project-Specific Standards

1. **Package Organization**
   - `cmd/`: Cobra CLI commands
   - `internal/`: Internal packages (not importable)
   - `internal/converter/`: Converter interfaces and registry
   - `internal/processors/`: Converter implementations
   - `internal/executor/`: FFmpeg and worker pool
   - `internal/config/`: Configuration management
   - `internal/ui/`: User interface (progress, output)

2. **Naming Conventions**
   - Use descriptive names
   - Converters: `{format}_to_{format}` (e.g., `mov_to_mp4`)
   - Interfaces: Noun (e.g., `Converter`)
   - Methods: Verb or VerbNoun (e.g., `Convert`, `ValidateInput`)

3. **Error Handling**
   - Return errors, don't panic
   - Wrap errors with context: `fmt.Errorf("context: %w", err)`
   - Log errors at appropriate levels

4. **Comments**
   - Public functions must have doc comments
   - Use `//` for single-line comments
   - Start comments with the name of the thing being described

### Code Quality

- **Linting**: All code must pass `golangci-lint`
- **Formatting**: Run `gofmt -s -w .` before committing
- **Dependencies**: Minimize external dependencies
- **Security**: No hardcoded credentials or secrets

## Testing

### Writing Tests

```go
// Example test
func TestMP4Converter_Validate(t *testing.T) {
    conv := NewMP4Converter()

    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid MOV file", "test.mov", false},
        {"invalid extension", "test.txt", true},
        {"non-existent file", "missing.mov", true},
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

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test ./internal/processors/mov_to_mp4 -v

# Run with race detector
go test -race ./...
```

### Test Requirements

- All new features must include tests
- Bug fixes should include regression tests
- Aim for >80% code coverage for new code
- Tests must pass before PR merge

## Pull Request Process

### Before Submitting

1. **Update from upstream**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes**
   - Follow coding standards
   - Add tests
   - Update documentation

4. **Run checks**
   ```bash
   make test
   golangci-lint run ./...
   go mod tidy
   ```

5. **Commit your changes**
   ```bash
   git add .
   git commit -m "type: description"
   ```

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): description

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `test`: Adding/updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(converter): add HEIC to JPG converter
fix(mp4): handle special characters in filenames
docs(readme): update installation instructions
test(executor): add ffmpeg wrapper tests
```

### Pull Request Template

1. Fill out the PR template completely
2. Reference related issues: `Fixes #123` or `Relates to #456`
3. Describe your changes clearly
4. List any breaking changes
5. Include testing steps

### Review Process

1. Automated checks must pass (CI, linting, tests)
2. Code review by maintainer(s)
3. Address review comments
4. Maintainer approval required for merge
5. Squash and merge strategy used

### After Merge

- Delete your feature branch
- Update your local repository
- Close related issues (if not auto-closed)

## Adding New Converters

See [docs/converters.md](docs/converters.md) for detailed guide.

### Quick Start

1. **Create processor directory**
   ```bash
   mkdir -p internal/processors/{format}_to_{format}
   ```

2. **Implement converter interface**
   ```go
   type MyConverter struct {
       // fields
   }

   func (c *MyConverter) Name() string { ... }
   func (c *MyConverter) Description() string { ... }
   func (c *MyConverter) SupportedInputs() []string { ... }
   func (c *MyConverter) OutputExtension() string { ... }
   func (c *MyConverter) Validate(input string) error { ... }
   func (c *MyConverter) Convert(input string, opts Options) (*Result, error) { ... }
   func (c *MyConverter) ConvertBatch(inputs []string, opts Options) ([]*Result, error) { ... }
   ```

3. **Register converter**
   ```go
   func init() {
       converter.Register(NewMyConverter())
   }
   ```

4. **Create command**
   ```bash
   # Create cmd/formats/{format}.go
   ```

5. **Add tests**
   ```bash
   # Create tests for your converter
   ```

6. **Update documentation**
   - Update README.md with new converter
   - Update CHANGELOG.md

## Questions?

- **General questions**: Use [GitHub Discussions](https://github.com/dusk-labs/sb/discussions)
- **Bug reports**: Create an [Issue](https://github.com/dusk-labs/sb/issues)
- **Feature requests**: Create an [Issue](https://github.com/dusk-labs/sb/issues)
- **Security issues**: See [SECURITY.md](SECURITY.md)

## License

By contributing to SB, you agree that your contributions will be licensed under the same license as the project (see repository root for license information).

---

Thank you for contributing to SB! 🎉
