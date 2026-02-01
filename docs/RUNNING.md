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

## Graceful Shutdown

The server handles `SIGINT` (Ctrl+C) and `SIGTERM` signals gracefully:

* Press `Ctrl+C` in the terminal, or
* Send `SIGTERM` to the process (e.g., `kill <pid>`)

The server will complete in-flight requests and shut down cleanly.

## Current Stage Limitations

**Stage 1** provides:
* HTTP server with `/healthcheck` endpoint
* Basic configuration via flags and environment variables
* Graceful shutdown handling

**Not yet available:**
* `/metrics` endpoint (Prometheus metrics)
* System metrics collection
* Docker/containerization support
