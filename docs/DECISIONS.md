# Architecture & Design Decisions

This document records **key architectural and design decisions** made during the development of **MyGoMetrics**.

The goal is not to present a "perfect" exporter, but to **make trade-offs explicit**, document intent, and show how the system evolved incrementally across releases.

---

## 1. Project Scope Definition

**Decision:** Build a host-level Prometheus exporter focused on system and Go runtime metrics only.

**Rationale:**

* Prometheus exporters are typically small, focused components
* Host-level metrics are broadly useful and easy to reason about
* Avoids becoming a general-purpose monitoring agent

**Non-Goals:**

* Full observability agent
* Log collection or tracing
* Application-specific business metrics

**Implications:**

* MyGoMetrics complements existing exporters rather than replacing them
* Feature growth is intentionally constrained

---

## 2. Collector-Based Architecture

**Decision:** Implement metrics collection via independent collectors.

**Rationale:**

* Encourages separation of concerns
* Allows collectors to be enabled or disabled independently
* Simplifies testing and future extensions

**Alternatives Considered:**

* Single monolithic collector
  → Rejected due to poor testability and coupling
* Deep abstraction hierarchy
  → Rejected as unnecessary for project size

**Implications:**

* Each collector owns its data source and error handling
* Failures are isolated to individual collectors

---

## 3. Prometheus-Agnostic Collectors

**Decision:** Collectors do not depend on Prometheus client types.

**Rationale:**

* Keeps data collection decoupled from exposition
* Makes collectors reusable and easier to test
* Prevents Prometheus-specific concerns from leaking into core logic

**Alternatives Considered:**

* Writing metrics directly using Prometheus primitives
  → Rejected due to tight coupling and reduced flexibility

**Implications:**

* Mapping to Prometheus metrics happens in a dedicated exporter layer
* Transport or exposition format can evolve independently

---

## 4. Metrics Source Selection

**Decision:** Use Go standard library and `gopsutil` for system metrics.

**Rationale:**

* `runtime` and `runtime/metrics` provide authoritative Go runtime data
* `gopsutil` is a mature, widely used library for system metrics
* Avoids platform-specific system calls where possible

**Alternatives Considered:**

* Writing OS-specific implementations
  → Rejected due to complexity and portability concerns

**Implications:**

* Metrics availability depends on OS support
* Behavior is consistent with other real-world exporters

---

## 5. Error Handling Strategy

**Decision:** Collector failures are logged and surfaced as partial metric loss, not fatal errors.

**Rationale:**

* Metrics exporters should favor availability
* A single failing metric should not break the entire `/metrics` endpoint
* Matches Prometheus ecosystem expectations

**Alternatives Considered:**

* Failing the scrape on any error
  → Rejected as too brittle

**Implications:**

* Metrics may be incomplete during transient failures
* Errors are observable via logs and exporter-level metrics

---

## 6. Configuration Strategy

**Decision:** Configuration via flags and environment variables, with explicit precedence.

**Rationale:**

* Simple and transparent configuration model
* Container- and Kubernetes-friendly
* Easy to document and reason about

**Configuration Precedence:**

1. Command-line flags
2. Environment variables
3. Defaults

**Alternatives Considered:**

* Configuration files
  → Deferred as unnecessary complexity

**Implications:**

* Configuration is static at startup
* Misconfiguration fails fast

---

## 7. HTTP Server Responsibilities

**Decision:** Expose only minimal HTTP endpoints.

**Implemented Endpoints:**

* `/metrics` — Prometheus scrape endpoint
* `/healthcheck` — basic liveness check

**Rationale:**

* Exporters should not expose rich APIs
* Smaller surface area reduces maintenance and risk

**Non-Goals:**

* Admin APIs
* Runtime reconfiguration

---

## 8. Concurrency Model

**Decision:** Periodic collection with a configurable interval.

**Rationale:**

* Predictable CPU and memory usage
* Aligns with Prometheus pull-based model
* Easier to reason about than event-driven collection

**Alternatives Considered:**

* Collect-on-scrape
  → Rejected due to variable scrape load
* Per-metric goroutines
  → Rejected as unnecessary complexity

**Implications:**

* Metrics may be slightly stale between intervals
* Resource usage remains bounded

---

## 9. Graceful Shutdown

**Decision:** Use `context.Context` for lifecycle management.

**Rationale:**

* Standard Go practice
* Ensures background collectors stop cleanly
* Avoids half-written metrics during shutdown

**Implications:**

* Shutdown waits for in-flight collection cycles
* Suitable for containerized environments

**Implementation:**

* Root context created in `main.go`
* Signal handling (SIGINT/SIGTERM) cancels context
* HTTP server performs graceful shutdown via `Shutdown(ctx)`

---

## 10. HTTP Router

**Decision:** Use the standard library `net/http` package for HTTP serving and routing.

**Rationale:**

* The server exposes minimal endpoints (`/healthcheck` and `/metrics`)
* No path parameters, complex routing, or middleware chains are required
* Using a third-party router would add an external dependency and conceptual overhead without tangible benefit

**Alternatives Considered:**

* go-chi/chi (lightweight router with path parameters, middleware, sub-routers)
  → Rejected as overcomplicating for minimal endpoints; net/http is sufficient

**Implications:**

* No router dependency
* If the HTTP surface grows significantly in later stages, revisiting Chi or a similar router remains an option

---

## 11. Testing Strategy

**Decision:** Test behavior and integration boundaries, not internal details.

**Approach:**

* Unit tests for collectors
* HTTP endpoint tests for `/metrics`
* Avoid excessive mocking of OS-level behavior

**Non-Goals:**

* Exhaustive platform simulation
* Kernel-level fault injection

**Implications:**

* Tests are stable and fast
* Some platform-specific behavior is validated manually

---

## 12. Containerization Choices

**Decision:** Use a minimal, non-root container image.

**Rationale:**

* Exporters are frequently deployed cluster-wide
* Smaller images reduce attack surface and startup time
* Non-root execution aligns with best practices

**Implementation Notes:**

* Multi-stage build
* Distroless or scratch-based runtime image

---

## Final Note

These decisions reflect **intentional trade-offs**, not missing knowledge.

MyGoMetrics is designed to be:

* simple
* predictable
* observable
* production-appropriate

Rather than exhaustive or over-engineered.
