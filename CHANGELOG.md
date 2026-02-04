# Changelog

All notable changes to the MyGoMetrics project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [0.7.0] - Stage 7: CI/CD

### Added

#### GitHub Actions CI/CD Workflow

- `**.github/workflows/ci.yml**`: Complete CI/CD pipeline for automated testing and releases
  - **Lint and Test Job**: Runs on every push/PR
    - `go vet ./...` for static code analysis
    - `go test ./...` for unit tests
    - Validates code quality on all relevant changes
  - **Docker Build Job**: Validates Docker image builds on all runs
    - Uses Docker Buildx for consistent builds
    - Context at repository root with `./Dockerfile`
    - Does not push on non-tag builds (validation only)
  - **Docker Push on Tag**: Automatic image publishing for releases
    - Triggered only on version tags (e.g., `v0.7.0`, `v*` pattern)
    - Builds and pushes to GitHub Container Registry (ghcr.io)
    - Uses `GITHUB_TOKEN` for authentication (no extra secrets required)
    - Single-architecture build (`linux/amd64`) with multi-arch support planned
    - Image naming: `ghcr.io/sheliakhin-golang-portfolio/mygometrics:<tag>`
  - **Path Filtering**: Workflow triggers only on relevant changes
    - Monitored paths: `.github/`, `cmd/`, `internal/`, `go.mod`, `Dockerfile`, `helm/`
    - Reduces unnecessary pipeline runs

#### Helm Chart Updates for Container Registry

- `**helm/mygometrics/Chart.yaml**`: Updated to version 0.7.0
  - Chart version incremented from 0.6.0 to 0.7.0
  - Reflects CI/CD integration and registry changes

- `**helm/mygometrics/values.yaml**`: Container registry configuration
  - `image.registry`: New field set to `ghcr.io` (GitHub Container Registry)
  - `image.repository`: Updated to `sheliakhin-golang-portfolio/mygometrics` (lowercase for compatibility)
  - `image.tag`: Defaults to `v0.7.0` for this release
  - Fully qualified image reference: `ghcr.io/sheliakhin-golang-portfolio/mygometrics:v0.7.0`

- `**helm/mygometrics/templates/deployment.yaml**`: Registry-aware image reference
  - Updated container image template to use both `image.registry` and `image.repository`
  - Previous format: `{{ .Values.image.repository }}:{{ .Values.image.tag }}`
  - New format: `{{ .Values.image.registry }}/{{ .Values.image.repository }}:{{ .Values.image.tag }}`
  - Supports full registry paths for container images

#### Documentation Updates

- `**docs/RUNNING.md**`: Added CI/CD and release section
  - How to cut a release using version tags
  - Image publishing location (ghcr.io)
  - Automated build and test process description
  - Tag format requirements (`v*` pattern)

- `**docs/DECISIONS.md**`: Section 14 (CI/CD and Release Strategy)
  - Rationale for automated CI/CD (reproducible builds, infrastructure component needs)
  - GitHub Actions implementation details
  - Container registry choice (ghcr.io)
  - Single-arch vs multi-arch build strategy
  - Implications for release process and versioning

### Technical Details

#### CI/CD Pipeline Architecture

- **Trigger Strategy**:
  - Push/PR: Runs linting, testing, and Docker build validation (no push)
  - Tag push (`v*`): Runs full pipeline including Docker push to registry
- **Go Version**: Automatically detected from `go.mod` (currently Go 1.25)
- **Working Directory**: Repository root (MyGoMetrics); all commands run from repo root
- **Docker Build Context**: Repository root (`.`), using `./Dockerfile`

#### Container Image Publishing

- **Registry**: GitHub Container Registry (ghcr.io)
  - No additional secrets required (uses `GITHUB_TOKEN`)
  - Public or private based on repository visibility
- **Image Naming Convention**: Lowercase for compatibility with registry requirements
  - Format: `ghcr.io/<owner>/<repository>:<tag>`
  - Example: `ghcr.io/sheliakhin-golang-portfolio/mygometrics:v0.7.0`
- **Tagging Strategy**:
  - Git tag becomes Docker image tag (e.g., `v0.7.0` → `:v0.7.0`)
  - Latest tag can be added for default branch tags (optional)

#### Multi-Architecture Support (Planned)

- Current implementation: Single-architecture (`linux/amd64`)
- Future enhancement: Multi-arch builds via Docker Buildx
  - Target platforms: `linux/amd64`, `linux/arm64`
  - Buildx driver: `docker-container` for multi-platform support
  - Push as manifest list for automatic platform selection
- **Feasibility**: Go 1.25 with `CGO_ENABLED=0` supports cross-compilation
- Current Dockerfile hardcodes `GOARCH=amd64`; multi-arch can be handled in CI without Dockerfile changes

#### Helm Chart Registry Integration

- Chart version and app version are decoupled
- Chart version increments with chart changes (0.7.0)
- Image registry/repository separation enables flexibility
- Deployment template constructs full image path from values
- Backward compatibility maintained via structured values

### Notes

- CI/CD workflow is triggered only on relevant path changes to optimize build times
- Docker images are validated on every build; published only on version tags
- GitHub Container Registry authentication uses built-in `GITHUB_TOKEN` (no manual secret setup)
- Single-architecture builds provide fast, reliable releases; multi-arch can be added incrementally
- Helm chart now supports any OCI-compatible registry via `image.registry` field
- Release process: tag commit with `v*` pattern, push tag, CI automatically builds and publishes
- All releases are reproducible via Git tags and automated builds

---

## [0.6.0] - Stage 6: Helm Chart

### Added

#### Helm Chart for Kubernetes Deployment

- `**helm/mygometrics/**`: Complete Helm chart for Kubernetes deployment
  - `Chart.yaml`: Chart metadata (version 0.6.0, appVersion 1.0.0)
  - `values.yaml`: Comprehensive configuration values for all app settings
  - `templates/deployment.yaml`: Kubernetes Deployment with health probes, Pod Security Standards (restricted profile)
  - `templates/service.yaml`: Kubernetes Service (ClusterIP by default)
  - `templates/servicemonitor.yaml`: Optional Prometheus Operator ServiceMonitor (disabled by default)
  - `templates/serviceaccount.yaml`: ServiceAccount creation (enabled by default)
  - `templates/_helpers.tpl`: Helm template helpers for labels and naming
  - `README.md`: Chart documentation with installation and configuration examples

#### Configuration via Helm Values

- All application configuration options available via Helm values:
  - `config.listenAddr` → `LISTEN_ADDR` environment variable
  - `config.collectInterval` → `COLLECT_INTERVAL` environment variable
  - `config.host` → `HOST` environment variable (metrics label)
  - `config.env` → `ENV` environment variable (metrics label)
  - `config.enabledCollectors` → `ENABLED_COLLECTORS` environment variable
  - `config.logLevel` → `LOG_LEVEL` environment variable
- Image configuration: `image.repository`, `image.tag`, `image.pullPolicy`
- Service configuration: `service.type`, `service.port`
- Resource requests/limits support via `resources` values
- ServiceAccount configuration via `serviceAccount` values

#### ServiceMonitor Support

- Optional Prometheus Operator integration:
  - `serviceMonitor.enabled`: Enable/disable ServiceMonitor creation
  - `serviceMonitor.interval`: Scrape interval configuration
  - `serviceMonitor.namespace`: Override ServiceMonitor namespace
  - `serviceMonitor.labels`: Labels for Prometheus Operator selector matching
- Automatic Prometheus discovery when ServiceMonitor is enabled and labels match

#### Documentation Updates

- `**docs/RUNNING.md**`: Added Kubernetes / Helm deployment section
  - Helm installation instructions (local chart, custom values, values file)
  - Configuration via Helm values
  - ServiceMonitor setup for Prometheus Operator
  - Upgrade and uninstallation commands
  - Reference to chart README and CONFIGURATION.md

### Technical Details

- Helm chart follows Kubernetes and Helm best practices:
  - Standard Helm labels (`app.kubernetes.io/name`, `app.kubernetes.io/instance`, etc.)
  - Template helpers for consistent naming and labeling
  - Conditional resource creation (ServiceMonitor, ServiceAccount)
  - Health probes (liveness and readiness) on `/healthcheck` endpoint
  - Pod Security Standards compliance (restricted profile): non-root user (65534), read-only root filesystem, no privilege escalation, all capabilities dropped, seccomp enabled
  - Resource requests/limits support for production deployments

- Chart versioning:
  - Chart version (0.6.0) increments with chart changes
  - App version (1.0.0) is fixed and decoupled from chart version
  - Chart version never decreases, even when syncing with app releases

- Kubernetes deployment features:
  - Single replica by default (configurable via `replicaCount`)
  - ClusterIP service by default (configurable via `service.type`)
  - Container port 9000 with named port `http`
  - Environment variables mapped from Helm values
  - Graceful shutdown support via Kubernetes signals

---

## [0.5.0] - Stage 5: Docker & Containerization

### Added

#### Docker Support

- `**Dockerfile**`: Production-ready multi-stage Docker build
  - Build stage: Uses `golang:1.25-alpine` for compilation
  - Static binary build with `CGO_ENABLED=0`, `GOOS=linux`, `GOARCH=amd64`
  - Binary optimization with `-ldflags="-w -s"` for smaller size
  - Runtime stage: Uses `gcr.io/distroless/static-debian12:nonroot` minimal image
  - Non-root user execution (distroless images use `nonroot` user by default)
  - Exposes port 9000 (default metrics and health check port)
  - Single binary entrypoint `/mygometrics`

- `**.dockerignore**`: Build context optimization
  - Excludes Git files, documentation, test files, IDE files, and build artifacts
  - Reduces Docker build context size for faster builds
  - Prevents accidental inclusion of sensitive files (`.env` files)

#### Documentation Updates

- `**docs/RUNNING.md**`: Added Docker usage section
  - Docker build instructions
  - Docker run examples with port mapping
  - Configuration via environment variables in containers
  - Notes on non-root execution and port exposure

- `**docs/CONFIGURATION.md**`: Added Docker configuration note
  - Clarifies that all existing flags and environment variables work in containers
  - Configuration precedence remains the same (flags > env vars > defaults)

### Technical Details

- Docker image follows containerization best practices:
  - Multi-stage build reduces final image size
  - Static binary enables use of minimal runtime images
  - Non-root execution improves security posture
  - Distroless base image eliminates shell and unnecessary tools
  - Build context optimization via `.dockerignore` improves build performance

- Container configuration:
  - Default listen address `:9000` works correctly in containers
  - All configuration options (flags, env vars) remain functional
  - Graceful shutdown handling works with container orchestration signals

---

## [0.4.0] - Stage 4: Control & Testing

### Added

#### Collector Control

- `**internal/config/config.go**`: Collector enable/disable configuration
  - `EnabledCollectors` field: slice of collector names to enable (default: empty = all enabled)
  - Flag: `-enabled-collectors` for comma-separated list of collector names
  - Environment variable: `ENABLED_COLLECTORS` for comma-separated list of collector names
  - Valid collector names: `cpu`, `memory`, `disk`, `runtime`
  - `parseCollectors()` helper function for parsing comma-separated collector lists with whitespace trimming
  - Updated GoDoc with collector configuration details

- `**cmd/main.go**`: Collector filtering logic
  - `filterCollectors()` function for filtering collectors based on enabled list
  - O(1) lookup using map for performance
  - Unknown collector names are ignored (no-op, no error)
  - Logging of enabled collectors at startup
  - If `EnabledCollectors` is empty, all collectors are enabled

#### Structured Logging

- `**internal/logger/logger.go**`: Structured logging with uber-go/zap
  - Global `Logger` variable for application-wide logging
  - `Init(level string)` function to initialize logger with specified level
  - `Sync()` function to flush buffered logs
  - Supports log levels: `debug`, `info`, `warn`, `error`
  - Production-ready JSON output format
  - Full GoDoc documentation

- `**internal/config/config.go**`: Log level configuration
  - `LogLevel` field: configurable logging level (default: `info`)
  - Flag: `-log-level` for setting log level
  - Environment variable: `LOG_LEVEL` for setting log level
  - Valid values: `debug`, `info`, `warn`, `error`

- `**cmd/main.go**`: Integrated structured logging
  - Logger initialization on startup with configured log level
  - Logger sync (flush) on application shutdown
  - Replaced standard library `log` with zap logger throughout application
  - Structured logging for startup, shutdown, signal handling, and error conditions
  - Fallback to standard library `log` for pre-initialization errors

- `**internal/exporter/registry.go**`: Logger integration in exporter
  - `NewRegistry()` now accepts `logger *zap.Logger` parameter
  - Structured logging for collector errors during metric updates
  - Error logs include collector name and error details as structured fields

#### Comprehensive Testing

- `**internal/exporter/prometheus_test.go**`: Exporter package tests
  - Table-driven tests with multiple scenarios
  - Mock collector implementation (`mockCollector`) for testing
  - Test cases:
    - All collectors succeed: verifies all metrics are exposed
    - All collectors fail: verifies HTTP handler still returns 200 (availability over completeness)
    - One failing collector: verifies partial results are exposed (resilience)
  - HTTP handler integration tests via `httptest`
  - No-op logger (`zap.NewNop()`) to avoid test output pollution
  - Tests verify error handling strategy (partial metric loss vs fatal errors)
  - Full test coverage of exporter resilience behavior

- `**docs/TESTING.md**`: Testing documentation
  - Quick start guide for running tests
  - Test coverage generation instructions
  - HTML coverage report generation
  - Test structure overview for collector and exporter packages
  - Testing strategy documentation (mock collectors only, no OS-level mocking)
  - Continuous integration pre-commit checks
  - Race detection instructions

#### Configuration Enhancements

- `**.env.example**`: Updated with new configuration variables
  - `ENABLED_COLLECTORS`: Comma-separated list of collector names to enable
  - `LOG_LEVEL`: Logging level configuration
  - Inline comments explaining valid values and defaults
  - Updated with Stage 4 configuration options

#### Dependencies

- Added `go.uber.org/zap` dependency
  - Industry-standard structured logging library for Go
  - Used for production-grade logging throughout the application
  - Provides JSON-formatted logs with structured fields
  - Zero-allocation logging performance

### Technical Details

- Collector filtering happens at startup; disabled collectors are not instantiated or run
- Empty or nil `EnabledCollectors` enables all collectors (default behavior)
- Unknown collector names in configuration are silently ignored (no-op)
- Structured logging uses JSON format for machine-readable logs
- Logger is initialized early in application startup
- Logger is synchronized (flushed) on graceful shutdown
- Test coverage focuses on behavior and integration boundaries, not internal details
- Exporter tests verify resilience: partial failures do not break `/metrics` endpoint
- Mock collectors are used exclusively in tests; no OS-level mocking

### Notes

- Collector enable/disable is useful for reducing resource usage or troubleshooting
- Structured logging improves observability in production environments
- Test coverage demonstrates production-grade error handling strategy
- One failing collector does not break the entire `/metrics` endpoint (availability over completeness)
- Tests verify that successful collectors' metrics are still exposed when others fail

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
