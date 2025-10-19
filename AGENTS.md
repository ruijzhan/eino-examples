# Repository Guidelines

## Project Structure & Module Organization

This repository contains examples for the Eino AI framework, organized into functional modules:

- **components/**: Core component examples (lambda, prompt, retriever, tool)
- **compose/**: Orchestration examples using Chain and Graph
- **flow/**: Flow-based programming patterns and stream processing
- **quickstart/**: Basic tutorials and getting-started guides
- **adk/**: Advanced development kit examples
- **internal/**: Shared utilities and helper functions
- **devops/**: Development and debugging tools

Source files use Go (`.go`) extension. Tests are colocated with source files using the `*_test.go` naming convention.

## Build, Test, and Development Commands

```bash
# Run all examples in a module
go run path/to/module/main.go

# Run tests
go test ./...

# Run single test
go test -run TestFunctionName ./path/to/package

# Run tests with verbose output
go test -v ./...

# Lint code
golangci-lint run

# Format code
gofumpt -w .
goimports -w .

# Build dependencies
go mod tidy
go mod vendor
```

## Coding Style & Naming Conventions

- **Indentation**: Use tabs (Go standard)
- **Naming**: Follow Go conventions - PascalCase for exported, camelCase for unexported
- **File structure**: Package-level documentation, imports, then implementation
- **Formatting**: Enforced by golangci-lint with gofumpt (extra-rules enabled) and goimports
- **Package structure**: Main executable in `main.go`, supporting logic in separate files
- **Imports**: Group standard library, third-party, then local imports (handled by goimports)
- **Error handling**: Use explicit error returns, avoid panic in production code
- **Types**: Use explicit types, avoid var declarations for complex types

## Testing Guidelines

- **Framework**: Standard Go testing package
- **Coverage**: Minimal test files present (`*_test.go`)
- **Naming**: Test functions follow `TestFunctionName` pattern
- **Location**: Tests colocated with implementation files
- **Running**: Use `go test ./...` to run all tests, `go test -run TestName ./pkg` for single test

## Commit & Pull Request Guidelines

### Commit Message Format
Follow conventional commits:
- `feat:` for new features
- `docs:` for documentation changes  
- `fix:` for bug fixes
- `refactor:` for code restructuring
- `chore:` for maintenance tasks

Examples from repository: `feat: complete ADK helloworld module`, `docs: update learning-path.md`

### Pull Request Requirements
- Use PR title format: `<type>(optional scope): <description>`
- Provide user-oriented descriptions for release notes
- Include Chinese translation of title when relevant
- Link related issues using `Fixes #<number>`
- Update user documentation if usage-level changes are required
- Fill out PR template completely

## Security & Configuration

- Store API keys in `.env` files (template: `.example.env`)
- Never commit secrets or credentials
- Use environment variables for configuration
- Follow security reporting guidelines in SECURITY section of README

## Linter Configuration
- Uses golangci-lint with gofumpt extra-rules enabled
- Disabled linters: errcheck, typecheck, staticcheck, unused, gosimple, ineffassign, gofumpt
- Enabled: goimports, gofmt
- Excludes mock files (*.mock.go) from analysis
