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
# Test all collectors
go test ./internal/collector

# Test exporter package
go test ./internal/exporter
```

Run specific tests:

```bash
# Run specific collector test
go test -v ./internal/collector -run TestRuntimeCollector

# Run specific exporter test
go test -v ./internal/exporter -run TestRegistry_Update_OneFailingCollectorDoesNotBreakMetrics
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

The project includes comprehensive unit tests organized by package:

### Collector Tests (`internal/collector/`)

- `runtime_test.go` - Go runtime metrics collector tests
- `memory_test.go` - System memory metrics collector tests
- `cpu_test.go` - CPU metrics collector tests
- `disk_test.go` - Disk metrics collector tests

**What's Tested:**

Each collector test verifies:

1. **Collector naming** - Returns expected collector name
2. **Metric collection** - Successfully collects metrics
3. **Metric validity** - Values are non-negative and correctly named
4. **Context handling** - Properly handles context cancellation
5. **Error handling** - Returns errors when context is cancelled

### Exporter Tests (`internal/exporter/`)

- `prometheus_test.go` - Prometheus registry and HTTP handler tests

**What's Tested:**

The exporter tests use **table-driven tests** with **mocked collectors** to verify:

1. **Resilience to collector failures** - One failing collector does not break `/metrics` endpoint
2. **Partial results** - Successful collectors' metrics are exposed even when others fail
3. **HTTP handler behavior** - `/metrics` endpoint returns HTTP 200 even with partial failures
4. **Metric presence** - Expected metrics from successful collectors appear in the response

**Testing Strategy:**

- **Mock collectors only** - No OS or system-level mocking
- **Test via HTTP handler** - Tests verify the real `/metrics` endpoint behavior
- **No-op logger in tests** - Uses `zap.NewNop()` to avoid test output pollution
- **Realistic scenarios** - Tests cover all success, all failure, and partial failure cases


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
