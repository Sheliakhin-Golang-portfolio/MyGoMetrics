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

## Environment File

An optional `.env` file can be used to set environment variables. See `.env.example` for the format.

**Note:** The `.env` file is loaded automatically if present in the current working directory. This is optional and primarily useful for local development.

---

## Validation

* The listen address cannot be empty
* Invalid addresses will cause the server to fail at startup
* All configuration is validated at startup before the server begins listening
