# Changelog

All notable changes to the MyGoMetrics project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [0.1.0] - Stage 1: Project initialization

### Added

#### Project Structure

- `**go.mod**`: Go module initialization
  - Module path: `github.com/Sheliakhin-Golang-portfolio/MyGoMetrics`
  - Go version: 1.25.0
  - Dependency: `github.com/joho/godotenv` for optional `.env` file loading

#### Directory Structure

- `**cmd/**`: Application entry point directory
- `**internal/**`: Private application code directory
- `**docs/**`: Documentation directory

#### Configuration Package

- `**internal/config/config.go**`: Configuration loading with flags and environment variables
  - `Config` struct with `ListenAddr string` field
  - `Load()` function that reads from flags, environment variables, and defaults
  - Precedence: flags > environment variables > defaults
  - Optional `.env` file loading via godotenv (ignores "file not found")
  - Flag: `-listen-addr` for HTTP server listen address (default: ":9000")
  - Environment variable: `LISTEN_ADDR` for HTTP server listen address
  - Validation: listen address cannot be empty
  - Full GoDoc documentation for exported types and `Load()`

#### HTTP Server Package

- `**internal/server/http.go**`: HTTP server with healthcheck and graceful shutdown
  - `Start(ctx context.Context, cfg config.Config) error` — starts server, blocks until context cancelled
  - Server listens on `cfg.ListenAddr` with configurable timeouts:
    - ReadTimeout: 15s, WriteTimeout: 15s, IdleTimeout: 60s
  - Route: `GET /healthcheck` — returns `200 OK` with `{"status":"ok"}` (JSON)
  - Method check: non-GET requests to `/healthcheck` return `405 Method Not Allowed`
  - Graceful shutdown: on context cancellation, `Shutdown()` with 10-second timeout
  - GoDoc documentation for `Start()` and `healthcheckHandler()`

#### Entry Point

- `**cmd/main.go**`: Application lifecycle and signal handling
  - Loads configuration via `config.Load()`; exits on error
  - Root context with cancellation for graceful shutdown
  - Signal handling for SIGINT and SIGTERM; cancels context on signal
  - Starts HTTP server via `server.Start(ctx, cfg)`; blocks until shutdown
  - Logging for startup, signal receipt, and shutdown completion

#### Configuration Template

- `**.env.example**`: Environment variable template
  - `LISTEN_ADDR` with example value `:9000`
  - Inline comments for HTTP server configuration and flag override

#### Documentation

- `**docs/CONFIGURATION.md**`: Configuration and environment variables
  - Configuration precedence (flags > env > defaults)
  - Server configuration table: variable, description, default, flag
  - Address format examples (`:9000`, `127.0.0.1:9000`, `0.0.0.0:9000`)
  - Usage examples (flag, env, `.env` file)
  - Validation rules and environment file notes
- `**docs/RUNNING.md**`: Build and run instructions
  - Prerequisites (Go version)
  - Build command (`go build ./cmd`)
  - Running with default and configured listen address
  - Configuration options (flags, env, `.env`)
  - Health check verification (`curl` example)
  - Graceful shutdown behavior (SIGINT/SIGTERM)
  - Current stage limitations (what is and is not yet available)
- `**docs/DECISIONS.md**`: Architecture Decision Records (ADRs)
  - Decision 1: Project scope (host-level Prometheus exporter, system and Go runtime metrics only)
  - Decision 2: Configuration strategy (flags and env with explicit precedence)
  - Decision 3: HTTP server responsibilities (minimal endpoints: `/healthcheck`, future `/metrics`)
  - Decision 4: Graceful shutdown using `context.Context`
  - Decision 5: HTTP router (standard library `net/http`, no third-party router for Stage 1)
  - Each decision includes rationale, alternatives considered, and implications

#### License

- `**LICENSE**`: Apache License 2.0
  - Full Apache License 2.0 text
  - Standard Apache 2.0 terms and conditions

#### Version Control

- `**.gitignore**`: Git ignore patterns for Go projects
  - Binaries, test artifacts, coverage, workspace files, `.env`

### Technical Details

- HTTP server exposes a single liveness endpoint (`/healthcheck`) for Stage 1
- Configuration is static at startup; misconfiguration fails fast
- All blocking operations respect context cancellation for graceful shutdown
- No Docker or Prometheus metrics in this stage; foundation for future exporter work
- Documentation follows consistent formatting and cross-references (CONFIGURATION, RUNNING, DECISIONS)
- Environment-driven configuration; no hardcoded listen address in code

---

## Notes

- All stages follow semantic versioning principles
- Each stage is designed to be runnable and maintainable
- Code follows Go 1.25+ conventions and idiomatic patterns
- Security and maintainability: validated config, explicit precedence, documented behavior
