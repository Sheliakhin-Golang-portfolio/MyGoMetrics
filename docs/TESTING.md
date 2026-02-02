# Testing Guide

## Quick Start

To run all tests:

```bash
go test ./...
```

This document describes how to run tests, generate test coverage, and understand the testing strategy for MyGoMetrics.

## Running Tests

Run all tests:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Run tests for a specific package:

```bash
go test ./internal/collector
```

Run specific tests:

```bash
go test -v ./internal/collector -run TestRuntimeCollector
```

Enable race detection:

```bash
go test -race ./...
```

## Test Coverage

Show coverage percentage:

```bash
go test -cover ./...
```

Generate coverage profile:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Generate HTML coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Test Structure

The project includes unit tests for all collectors in `internal/collector/`:

- `runtime_test.go` - Go runtime metrics collector tests
- `memory_test.go` - System memory metrics collector tests
- `cpu_test.go` - CPU metrics collector tests
- `disk_test.go` - Disk metrics collector tests

### What's Tested

Each collector test verifies:

1. **Collector naming** - Returns expected collector name
2. **Metric collection** - Successfully collects metrics
3. **Metric validity** - Values are non-negative and correctly named
4. **Context handling** - Properly handles context cancellation
5. **Error handling** - Returns errors when context is cancelled


## Continuous Integration

Pre-commit checks:

```bash
go test ./...
go test -race ./...
go vet ./...
go fmt ./...
```

CI pipeline command:

```bash
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
```
