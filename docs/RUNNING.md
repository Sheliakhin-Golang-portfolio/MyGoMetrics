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

### Configuration Options

You can configure the listen address using:

1. **Command-line flags** (highest precedence):
   ```bash
   ./cmd -listen-addr :8080
   ```

2. **Environment variables**:
   ```bash
   export LISTEN_ADDR=:8080
   ./cmd
   ```

3. **`.env` file** (optional):
   ```bash
   # Copy .env.example to .env and modify as needed
   cp .env.example .env
   ./cmd
   ```

See [CONFIGURATION.md](./CONFIGURATION.md) for details on configuration precedence.

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

**Stage 3** provides:
* HTTP server with `/healthcheck` endpoint
* `/metrics` endpoint for Prometheus scraping
* Basic configuration via flags and environment variables
* Metrics collection configuration (interval, host, env labels)
* Graceful shutdown handling
* Core metrics collection (CPU, memory, disk, runtime) via collector package
* Prometheus exporter with periodic collection and metric exposition

**Not yet available:**
* Docker/containerization support
* Helm charts for Kubernetes deployment
