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

**Stage 5 (v0.5.0)** provides:
* All Stage 4 features plus:
* **Docker containerization support** - Multi-stage Dockerfile with distroless runtime
* **Production-ready container image** - Non-root execution, minimal base image

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

---

## Running with Docker

### Building the Docker Image

Build the Docker image from the `MyGoMetrics` directory:

```bash
docker build -t mygometrics .
```

This creates a multi-stage build that:
- Compiles a static Go binary in the build stage
- Copies only the binary to a minimal distroless runtime image
- Results in a small, secure container image

### Running the Container

Run the container with default settings (listens on port 9000):

```bash
docker run --rm -p 9000:9000 mygometrics
```

**Configuration Options:**

All configuration options work the same way in containers:

```bash
# Using environment variables
docker run --rm -p 9000:9000 \
  -e LISTEN_ADDR=:9000 \
  -e COLLECT_INTERVAL=30s \
  -e ENV=production \
  -e LOG_LEVEL=info \
  mygometrics

# Using command-line flags
docker run --rm -p 9000:9000 \
  mygometrics -listen-addr :9000 -env production -log-level info

# Enable only specific collectors
docker run --rm -p 9000:9000 \
  -e ENABLED_COLLECTORS=cpu,memory \
  mygometrics
```

**Port Mapping:**

The container exposes port 9000 by default. Map it to a different host port if needed:

```bash
docker run --rm -p 8080:9000 mygometrics
```

**Security Notes:**

- The container runs as a non-root user (`nonroot`) by default
- The distroless base image contains no shell or unnecessary tools
- The binary is statically compiled with CGO disabled for maximum portability

**Verification:**

Once running, verify the container is healthy:

```bash
curl http://localhost:9000/healthcheck
curl http://localhost:9000/metrics
```

---

## Kubernetes / Helm Deployment

MyGoMetrics can be deployed to Kubernetes using the provided Helm chart.

### Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- Prometheus Operator (optional, for ServiceMonitor support)

### Installation

**Install from local chart:**

```bash
helm install mygometrics ./helm/mygometrics
```

**Install with custom values:**

```bash
helm install mygometrics ./helm/mygometrics \
  --set image.repository=your-registry/mygometrics \
  --set image.tag=0.6.0 \
  --set config.env=production \
  --set config.logLevel=info
```

**Install with values file:**

```bash
helm install mygometrics ./helm/mygometrics -f my-values.yaml
```

### Configuration

All application configuration options are available via Helm values under `config.*`:

- `config.listenAddr` → `LISTEN_ADDR` environment variable
- `config.collectInterval` → `COLLECT_INTERVAL` environment variable
- `config.host` → `HOST` environment variable (metrics label)
- `config.env` → `ENV` environment variable (metrics label)
- `config.enabledCollectors` → `ENABLED_COLLECTORS` environment variable
- `config.logLevel` → `LOG_LEVEL` environment variable

See `helm/mygometrics/values.yaml` for all available options and defaults. For detailed configuration semantics, see [CONFIGURATION.md](./CONFIGURATION.md).

### ServiceMonitor (Prometheus Operator)

To enable automatic Prometheus discovery, enable the ServiceMonitor:

```bash
helm install mygometrics ./helm/mygometrics \
  --set serviceMonitor.enabled=true \
  --set serviceMonitor.labels.release=prometheus
```

Or in `values.yaml`:

```yaml
serviceMonitor:
  enabled: true
  interval: "30s"
  labels:
    release: prometheus
```

**Note:** The `serviceMonitor.labels` must match your Prometheus instance's selector for automatic discovery.

### Upgrading

```bash
helm upgrade mygometrics ./helm/mygometrics
```

### Uninstallation

```bash
helm uninstall mygometrics
```

For more details, see [Helm`s Readme.md](../helm/mygometrics/README.md).