# Running MyGoMetrics

This document describes how to build and run MyGoMetrics.

## Prerequisites

* Go ≥ 1.25
* No external dependencies required for basic operation

## Building

Build the binary from the `MyGoMetrics` directory:

```bash
go build ./cmd
```

This creates a `cmd` executable (or `cmd.exe` on Windows) in the current directory.

## Running

### Basic Usage

Run the binary with default settings (listens on `:9000`):

```bash
./cmd
```

### Configuration

The application can be configured using command-line flags, environment variables, or a `.env` file.

**Quick Examples:**

```bash
# Using command-line flags (highest precedence)
./cmd -listen-addr :8080 -collect-interval 30s -env production

# Using environment variables
export LISTEN_ADDR=:8080
export COLLECT_INTERVAL=30s
export ENABLED_COLLECTORS=cpu,memory
./cmd

# Using .env file (copy .env.example to .env and modify)
cp .env.example .env
./cmd
```

**Common Configuration Examples:**

```bash
# Enable only specific collectors
./cmd -enabled-collectors cpu,memory

# Change log level for debugging
./cmd -log-level debug

# Production configuration
./cmd -listen-addr :9000 -env production -log-level info
```

For a complete list of all configuration options, flags, environment variables, defaults, and configuration precedence rules, see [CONFIGURATION.md](./CONFIGURATION.md).

## Health Check

Once running, verify the server is healthy:

```bash
curl http://localhost:9000/healthcheck
```

Expected response: `{"status":"ok"}` with HTTP 200 status.

## Metrics Endpoint

The `/metrics` endpoint exposes Prometheus-formatted metrics for scraping:

```bash
curl http://localhost:9000/metrics
```

This endpoint returns all collected metrics in Prometheus text format, including:

* **CPU metrics**: `mygometrics_cpu_usage_percent`
* **Memory metrics**: `mygometrics_memory_used_bytes`, `mygometrics_memory_total_bytes`
* **Disk metrics**: `mygometrics_disk_read_bytes`, `mygometrics_disk_write_bytes`
* **Runtime metrics**: `mygometrics_runtime_goroutines`, `mygometrics_runtime_gc_cycles`, `mygometrics_runtime_heap_alloc_bytes`

All metrics include `host` and `env` labels as configured.

**Prometheus Configuration:**

To scrape metrics from MyGoMetrics, add a scrape config to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'mygometrics'
    static_configs:
      - targets: ['localhost:9000']
```

Metrics are collected periodically (default: every 15 seconds) and updated in the Prometheus registry. The `/metrics` endpoint serves the current state of all metrics.

## Graceful Shutdown

The server handles `SIGINT` (Ctrl+C) and `SIGTERM` signals gracefully:

* Press `Ctrl+C` in the terminal, or
* Send `SIGTERM` to the process (e.g., `kill <pid>`)

The server will complete in-flight requests and shut down cleanly.

## Current Stage Features

**Stage 4 (v0.4.0)** provides:
* HTTP server with `/healthcheck` endpoint
* `/metrics` endpoint for Prometheus scraping
* Full configuration via flags, environment variables, and `.env` files
* Metrics collection configuration (interval, host, env labels)
* **Collector enable/disable by name** - selectively enable collectors
* **Structured logging with uber-go/zap** - JSON logs with configurable levels
* Graceful shutdown handling
* Core metrics collection (CPU, memory, disk, runtime) via collector package
* Prometheus exporter with periodic collection and metric exposition
* **Production-grade error handling** - one failing collector does not break `/metrics`
* **Comprehensive test coverage** - including exporter tests with mocked collectors

**Not yet available:**
* Docker/containerization support
* Helm charts for Kubernetes deployment
