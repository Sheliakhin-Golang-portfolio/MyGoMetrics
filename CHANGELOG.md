# Changelog

All notable changes to the MyGoMetrics project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [0.3.0] - Stage 3: Prometheus Exporter

### Added

#### Exporter Package

- `**internal/exporter/prometheus.go**`: Prometheus metric definitions and mapping
  - Metric type definitions (Gauge vs Counter)
  - Mapping from collector snapshot names to Prometheus metrics
  - Support for Gauge metrics (CPU usage, memory, runtime goroutines/heap)
  - Support for Counter metrics (disk I/O, GC cycles) with delta tracking
  - Label support (`host`, `env`) on all metrics
  - Full GoDoc documentation

- `**internal/exporter/registry.go**`: Prometheus registry bridge
  - `Registry` struct bridging collectors to Prometheus metrics
  - `NewRegistry()` constructor with collectors and label values
  - `Update(ctx)` method for periodic metric collection and updates
  - `Handler()` method returning HTTP handler for `/metrics` endpoint
  - Per-collector error handling (log and continue strategy)
  - Counter delta computation for cumulative metrics
  - Full GoDoc documentation

#### Configuration Enhancements

- `**internal/config/config.go**`: Extended configuration
  - `CollectInterval` field: configurable interval between collector runs (default: 15s)
  - `HostName` field: hostname label value for Prometheus metrics (default: `os.Hostname()`)
  - `Env` field: environment label value for Prometheus metrics (default: empty string)
  - Flags: `-collect-interval`, `-host`, `-env`
  - Environment variables: `COLLECT_INTERVAL`, `HOST`, `ENV`
  - Duration parsing and validation for collect interval
  - Updated GoDoc with new flags and environment variables

#### HTTP Server Integration

- `**internal/server/http.go**`: Metrics endpoint support
  - `Start()` function now accepts optional `metricsHandler http.Handler` parameter
  - `/metrics` endpoint registration when handler is provided
  - Backward compatible: handler can be `nil` (no `/metrics` endpoint)
  - Updated GoDoc to reflect metrics endpoint support

#### Main Application Wiring

- `**cmd/main.go**`: Complete exporter integration
  - Collector instantiation (CPU, Memory, Disk, Runtime)
  - Exporter registry creation with host/env labels
  - Periodic collection goroutine running at `CollectInterval`
  - Initial collection on startup
  - Metrics handler passed to HTTP server
  - Startup logging for collection interval

#### Prometheus Metrics

All metrics exposed with `mygometrics_` prefix and `host`/`env` labels:

- `mygometrics_cpu_usage_percent` (Gauge): CPU usage percentage (0-100)
- `mygometrics_memory_used_bytes` (Gauge): Memory used in bytes
- `mygometrics_memory_total_bytes` (Gauge): Total memory available in bytes
- `mygometrics_disk_read_bytes` (Counter): Total bytes read from disk (cumulative)
- `mygometrics_disk_write_bytes` (Counter): Total bytes written to disk (cumulative)
- `mygometrics_runtime_goroutines` (Gauge): Current number of goroutines
- `mygometrics_runtime_gc_cycles` (Counter): Total GC cycles since program start (cumulative)
- `mygometrics_runtime_heap_alloc_bytes` (Gauge): Bytes allocated on the heap

#### Dependencies

- Added `github.com/prometheus/client_golang` dependency
  - Used for Prometheus metric types (Gauge, Counter)
  - Used for Prometheus registry and HTTP handler
  - Official Prometheus Go client library

#### Documentation

- `**.env.example**`: Updated with new configuration variables
  - `COLLECT_INTERVAL`: Interval between collector runs
  - `HOST`: Hostname label value
  - `ENV`: Environment label value
  - Inline comments explaining defaults and usage

- `**docs/CONFIGURATION.md**`: Metrics collection configuration section
  - New configuration table for metrics collection settings
  - Examples for interval, host, and env configuration
  - Validation rules for collect interval

- `**docs/RUNNING.md**`: Metrics endpoint documentation
  - `/metrics` endpoint usage and examples
  - Prometheus scrape configuration example
  - List of exposed metrics
  - Updated stage features section

- `**docs/DECISIONS.md**`: Architecture Decision Records updated for Stage 3
  - Decision 2: Collector-Based Architecture (independent collectors for separation of concerns)
  - Decision 3: Prometheus-Agnostic Collectors (decoupled data collection from exposition)
  - Decision 4: Metrics Source Selection (Go stdlib and gopsutil rationale)
  - Decision 5: Error Handling Strategy (partial metric loss vs fatal errors)
  - Decision 8: Concurrency Model (periodic collection with configurable interval)
  - Updated existing decisions to reflect Stage 3 implementation (HTTP endpoints, graceful shutdown)
  - Document now covers complete project evolution from Stage 1 through Stage 3

### Technical Details

- Periodic collection model: metrics collected at configurable intervals (default: 15s)
- Counter delta tracking: cumulative counters (disk I/O, GC cycles) track previous values and add deltas
- Error handling: per-collector failures are logged; partial metrics are still served
- Prometheus-agnostic collectors: collectors remain unchanged; no Prometheus imports in `internal/collector`
- Label support: all metrics include `host` and `env` labels for multi-instance deployments
- Graceful shutdown: periodic collection goroutine respects context cancellation

### Notes

- Metrics are now exposed via HTTP `/metrics` endpoint for Prometheus scraping
- Collection happens periodically in background; `/metrics` serves current registry state
- Counter metrics handle resets gracefully (negative deltas are ignored)
- Hostname is automatically detected; can be overridden via config
- Environment label is optional (empty string if not set)

---

## [0.2.0] - Stage 2: Core Metrics Collection

### Added

#### Collector Package

- `**internal/collector/collector.go**`: Core collector interface and types
  - `Snapshot` struct: Prometheus-agnostic metric representation with name, value, and optional labels
  - `Collector` interface: `Name() string` and `Collect(ctx context.Context) ([]Snapshot, error)`
  - Full GoDoc documentation for exported types and interface
  - Labels may be nil to indicate no labels (documented behavior)

#### CPU Collector

- `**internal/collector/cpu.go**`: CPU usage metrics collection
  - `CPUCollector` struct implementing `Collector` interface
  - `NewCPUCollector()` constructor function
  - Uses `gopsutil/v3/cpu.PercentWithContext` with 1-second interval
  - Metric: `cpu_usage_percent` (0-100 percentage)
  - Context-aware collection with proper error handling

#### Memory Collector

- `**internal/collector/memory.go**`: Memory metrics collection
  - `MemoryCollector` struct implementing `Collector` interface
  - `NewMemoryCollector()` constructor function
  - Uses `gopsutil/v3/mem.VirtualMemoryWithContext`
  - Metrics: `memory_used_bytes`, `memory_total_bytes`
  - Context-aware collection with proper error handling

#### Disk Collector

- `**internal/collector/disk.go**`: Disk I/O metrics collection
  - `DiskCollector` struct implementing `Collector` interface
  - `NewDiskCollector()` constructor function
  - Uses `gopsutil/v3/disk.IOCountersWithContext`
  - Aggregates I/O counters across all devices
  - Metrics: `disk_read_bytes`, `disk_write_bytes`
  - Handles platform-specific limitations gracefully

#### Runtime Collector

- `**internal/collector/runtime.go**`: Go runtime metrics collection
  - `RuntimeCollector` struct implementing `Collector` interface
  - `NewRuntimeCollector()` constructor function
  - Uses standard library `runtime` package (no external dependencies)
  - Uses `runtime.ReadMemStats` for compatibility (runtime/metrics noted for future)
  - Metrics: `runtime_goroutines`, `runtime_gc_cycles`, `runtime_heap_alloc_bytes`
  - Context-aware collection

#### Testing

- `**internal/collector/*_test.go**`: Comprehensive unit tests for all collectors
  - Table-driven test patterns following Go best practices
  - Tests for `Name()` method on all collectors
  - Tests for `Collect()` method with valid contexts
  - Tests for context cancellation behavior
  - Test coverage: 90.2% of statements
  - `test_helpers.go`: Shared test helper functions

#### Dependencies

- Added `github.com/shirou/gopsutil/v3` dependency for system metrics
  - Used for CPU, memory, and disk metrics collection
  - Industry-standard library used by real exporters

### Technical Details

- All collectors follow the Prometheus-agnostic design (no Prometheus client types)
- Error handling strategy: collectors return partial results with errors rather than failing completely
- Context support: all collectors respect context cancellation and timeouts
- Metrics are collected in-memory but not yet exposed via HTTP endpoints
- Collectors are testable independently without requiring a running server
- Full GoDoc documentation for all exported types and functions

### Notes

- Metrics collection is functional but not yet exposed via `/metrics` endpoint (Stage 3)
- No Prometheus client integration in this stage (per design decisions)
- Disk I/O metrics may not be available on all platforms (handled gracefully)
- Runtime collector uses `runtime.ReadMemStats`; `runtime/metrics` noted for future consideration

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
