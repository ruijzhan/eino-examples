---
trigger: always_on
---

# Eino Examples Development Guide

## Project Overview

This repository contains comprehensive examples and demonstrations for the **Eino framework** - a Go-based AI orchestration framework developed by CloudWeGo. The project provides practical implementations to help developers understand and utilize Eino's features for building AI applications.

**Key Characteristics:**
- **Language**: Go (1.24+ recommended)
- **Framework**: Eino v0.5.0 (CloudWeGo's AI orchestration framework)
- **License**: Apache 2.0
- **Organization**: Examples are structured by complexity and use case

## Technology Stack & Dependencies

### Core Dependencies
- **Eino Framework**: `github.com/cloudwego/eino v0.5.0`
- **Extension Components**: Various `github.com/cloudwego/eino-ext/*` packages
- **Model Providers**: OpenAI, DeepSeek, Ollama, Ark, and others
- **Infrastructure**: Redis, Volcengine, Coze Loop integration
- **Web Framework**: Hertz (CloudWeGo's HTTP framework)
- **Observability**: OpenTelemetry, Langfuse, APM+ integration

### Development Tools
- **Go Version**: 1.24.0+ (toolchain go1.24.4)
- **Linting**: golangci-lint with custom configuration
- **Testing**: Standard Go testing with benchmark support
- **CI/CD**: GitHub Actions

## Project Structure

```
eino-examples/
├── adk/                    # Agent Development Kit examples
│   ├── helloworld/         # Basic agent examples
│   ├── intro/             # Introduction to agent concepts
│   ├── multiagent/        # Multi-agent systems
│   └── common/            # Shared utilities
├── components/            # Component usage examples
│   ├── document/          # Document processing (parsers, loaders)
│   ├── lambda/            # Lambda function examples
│   ├── prompt/            # Prompt engineering
│   ├── retriever/         # Data retrieval systems
│   └── tool/              # Tool integration examples
├── compose/               # Orchestration examples
│   ├── chain/             # Chain orchestration
│   ├── graph/             # Graph-based workflows
│   └── workflow/          # Workflow management
├── flow/                  # Flow-based programming
│   ├── chat/              # Chat flow examples
│   ├── eino_assistant/    # Assistant implementations
│   └── todoagent/         # Task-oriented agents
├── quickstart/            # Quick start examples
│   ├── agent/             # Basic agent patterns
│   └── chat/              # Simple chat implementations
├── devops/                # DevOps and monitoring
└── internal/              # Internal utilities
```

## Development Commands

### Basic Development
```bash
# Run any example (navigate to directory first)
cd quickstart/chat
go run main.go

# Build specific examples
go build -o bin/chat quickstart/chat/main.go
```

### Testing
```bash
# Run all tests with coverage
go test -race -covermode=atomic -coverprofile=coverage.out ./...

# Run benchmarks
go test -bench=. -benchmem -run=none ./...

# Test specific package
go test ./components/document/parser/...
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run linter (uses .golangci.yaml config)
golangci-lint run --timeout 5m

# Check imports
goimports -w .
```

### Multi-module Management
The project contains multiple Go modules in subdirectories:
```bash
# For quickstart examples
cd quickstart/eino_assistant
go mod tidy
go run main.go

# For flow examples
cd flow/agent/deer-go
go mod tidy
go run main.go
```

## Architecture Overview

### Eino Framework Concepts
1. **Components**: Reusable building blocks (models, retrievers, tools)
2. **Orchestration**: Chain, Graph, and Workflow composition patterns
3. **Flows**: Stream-based data processing pipelines
4. **Agents**: Autonomous AI entities with tools and memory

### Key Patterns Demonstrated
- **Agent Development**: From simple chatbots to complex multi-agent systems
- **Document Processing**: Various parsers (HTML, PDF, text) and loaders
- **Model Integration**: Multiple LLM providers (OpenAI, DeepSeek, Ollama)
- **Retrieval Systems**: Vector databases, search, and knowledge bases
- **Tool Integration**: External APIs and function calling
- **Workflow Orchestration**: Complex AI pipeline construction

## Development Guidelines

### Code Style
- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Adhere to Apache 2.0 license headers (automatically checked in CI)
- Commit messages follow [AngularJS conventions](https://docs.google.com/document/d/1QrDFcIiPjSLDn3EL15IJygNPiHORgU1_OOAqWjiDU5Y/edit)

### Code Investigation
- When investigating the implementation of referenced libraries, locate and examine relevant code within the `vendor` directory
- **Do not modify any files** in the vendor directory during investigation
- The vendor directory contains vendored dependencies for reproducible builds

### Contributing
1. Fork the repository
2. Create feature branch from `develop`
3. Include tests for new functionality
4. Run linting and tests locally
5. Submit PR to `develop` branch

### CI/CD Pipeline
- **PR Checks**: License headers, spell check, linting
- **Tests**: Unit tests with coverage, benchmarks
- **Quality Gates**: Automatic linting with GitHub Actions

## Getting Started

1. **Prerequisites**: Go 1.24+, API keys for model providers
2. **Quick Start**: `quickstart/` directory contains the simplest examples
3. **Progressive Learning**:
   - Start with `quickstart/chat` for basic LLM usage
   - Move to `components/` for individual feature exploration
   - Explore `compose/` for orchestration patterns
   - Build complex systems with `adk/` and `flow/`

## Example Categories

### By Complexity
- **Beginner**: `quickstart/` - Basic usage patterns
- **Intermediate**: `components/` - Individual feature exploration
- **Advanced**: `compose/` - Orchestration and workflow design
- **Expert**: `adk/` - Multi-agent systems and complex applications

### By Use Case
- **Chat Applications**: Basic chat to complex assistants
- **Document Processing**: PDF, HTML parsing and transformation
- **Knowledge Bases**: Retrieval and RAG implementations
- **Workflow Automation**: Business process automation
- **Multi-Agent Systems**: Collaborative AI agents

## Environment Setup

Most examples require environment variables for API keys. A `.env` file is provided at the workspace root with these configurations:

```bash
# Load environment variables from the .env file
source .env

# The main .env file contains:
# OPENAI_API_KEY=your-api-key
# OPENAI_MODEL_NAME=glm-4.5
# OPENAI_BASE_URL=https://api.z.ai/api/coding/paas/v4
```

Some examples may have additional environment variables or their own `.env` files. Check individual example directories for specific setup requirements.