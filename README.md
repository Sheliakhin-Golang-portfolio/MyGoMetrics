# MyGoMetrics

A production-style Prometheus exporter written in Go that collects host-level system metrics and Go runtime metrics, exposing them via a standard `/metrics` endpoint for Prometheus scraping.

---

## What This Project Demonstrates

This project was created as a portfolio application with a focus on backend engineering practices for observability and metrics collection.

It demonstrates:

- Building stateless Prometheus exporters in Go
- System metrics collection using industry-standard libraries
- Go runtime metrics via `runtime/metrics`
- Prometheus client library integration
- Graceful shutdown and lifecycle management
- Configuration via flags and environment variables
- Production-oriented error handling and observability

---

## Product Overview

MyGoMetrics is a Prometheus exporter designed to collect and expose host-level system metrics (CPU, memory, disk) and Go runtime metrics.

The service architecture follows a clear flow:

- **Collectors** gather metrics from various sources (system, Go runtime)
- **Prometheus Registry** aggregates collected metrics
- **HTTP Server** exposes `/metrics` endpoint for Prometheus scraping
- **Health Check** provides basic liveness endpoint

The project is designed to demonstrate how real Go exporters are built and evolved, focusing on clean architecture, reliability, and operational visibility — the core concerns of production monitoring infrastructure.

**Current Stage:** Stage 5 - Docker & Containerization

---

## Architecture Overview

At a high level:

- **Collectors** independently gather metrics from system and runtime sources
- **Prometheus Registry** aggregates metrics from all collectors
- **HTTP Server** serves metrics on `/metrics` endpoint
- **Configuration** via flags and environment variables with explicit precedence

The service is **stateless** — no internal database or persistence layer. All metrics are collected on-demand or periodically and exposed via HTTP.

Detailed architectural decisions are documented separately.

---

## Documentation

All detailed documentation is moved to `docs/` folder:

- [Architectural and design decisions](docs/DECISIONS.md)
- [Configuration & Environment Variables](docs/CONFIGURATION.md)
- [Running Locally](docs/RUNNING.md)
- [Testing Guide](docs/TESTING.md)

---

## CHANGELOG

As there are no multiple releases during development, all version changes are contained in [CHANGELOG.md](CHANGELOG.md) file.

---

## License

This project is developed for educational and portfolio purposes.
Licensed under the Apache License 2.0.
