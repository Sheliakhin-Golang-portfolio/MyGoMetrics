# Configuration & Environment Variables

This document describes all configuration options and environment variables used by the MyGoMetrics service.

---

## Configuration Precedence

Configuration values are resolved in the following order (highest to lowest precedence):

1. **Command-line flags**
2. **Environment variables**
3. **Default values**

If a value is set via multiple methods, the highest precedence source takes effect.

---

## Server Configuration

| Variable | Description | Default | Flag |
|----------|-------------|---------|------|
| `LISTEN_ADDR` | HTTP server listen address and port | `:9000` | `-listen-addr` |

Used by:
- MyGoMetrics HTTP server

**Address Format:**
- `:9000` - Listen on all interfaces, port 9000
- `127.0.0.1:9000` - Listen on localhost only, port 9000
- `0.0.0.0:9000` - Listen on all interfaces, port 9000 (explicit)

**Examples:**

```bash
# Using command-line flag
go run . -listen-addr :8080

# Using environment variable
export LISTEN_ADDR=:8080

# Using .env file (loaded automatically if present)
# See .env.example for format
```

---

## Metrics Collection Configuration

| Variable | Description | Default | Flag |
|----------|-------------|---------|------|
| `COLLECT_INTERVAL` | Interval between collector runs | `15s` | `-collect-interval` |
| `HOST` | Hostname label value for Prometheus metrics | `os.Hostname()` or `"unknown"` | `-host` |
| `ENV` | Environment label value for Prometheus metrics | `""` (empty) | `-env` |
| `ENABLED_COLLECTORS` | Comma-separated list of collector names to enable | `""` (empty = all enabled) | `-enabled-collectors` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` | `-log-level` |

Used by:
- Prometheus exporter (periodic metrics collection)

**Interval Format:**
- Duration strings: `15s`, `30s`, `1m`, `2m30s`
- Must be a positive duration

**Collector Names:**
- Valid collector names: `cpu`, `memory`, `disk`, `runtime`
- If `ENABLED_COLLECTORS` is empty or not set, all collectors are enabled
- Unknown collector names are ignored (no-op)
- Collector names are case-sensitive

**Examples:**

```bash
# Using command-line flags
go run . -collect-interval 30s -host myserver -env production

# Enable only CPU and memory collectors
go run . -enabled-collectors cpu,memory

# Using environment variables
export COLLECT_INTERVAL=30s
export HOST=myserver
export ENV=production
export ENABLED_COLLECTORS=cpu,memory,disk
export LOG_LEVEL=info
go run .

# Using .env file
# See .env.example for format
```

---

## Environment File

An optional `.env` file can be used to set environment variables. See `.env.example` for the format.

**Note:** The `.env` file is loaded automatically if present in the current working directory. This is optional and primarily useful for local development.

---

## Validation

* The listen address cannot be empty
* Invalid addresses will cause the server to fail at startup
* Collect interval must be a positive duration
* Invalid collect interval format will cause the server to fail at startup
* All configuration is validated at startup before the server begins listening
