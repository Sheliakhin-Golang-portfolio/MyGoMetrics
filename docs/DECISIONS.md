# Architecture & Design Decisions (Stage 1)

This document records **key architectural and design decisions** relevant to **Stage 1** of MyGoMetrics development.

**Note:** This file will be updated as development progresses and new decisions are made.

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

## 2. Configuration Strategy

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

## 3. HTTP Server Responsibilities

**Decision:** Expose only minimal HTTP endpoints.

**Implemented Endpoints:**

* `/healthcheck` — basic liveness check (Stage 1)
* `/metrics` — Prometheus scrape endpoint (future stage)

**Rationale:**

* Exporters should not expose rich APIs
* Smaller surface area reduces maintenance and risk

**Non-Goals:**

* Admin APIs
* Runtime reconfiguration

---

## 4. Graceful Shutdown

**Decision:** Use `context.Context` for lifecycle management.

**Rationale:**

* Standard Go practice
* Ensures background collectors stop cleanly
* Avoids half-written metrics during shutdown

**Implications:**

* Shutdown waits for in-flight collection cycles
* Suitable for containerized environments

**Stage 1 Implementation:**

* Root context created in `main.go`
* Signal handling (SIGINT/SIGTERM) cancels context
* HTTP server performs graceful shutdown via `Shutdown(ctx)`

---

## 5. HTTP Router

**Decision:** Use the standard library `net/http` package for HTTP serving and routing.

**Rationale:**

* The Stage 1 server exposes a single endpoint (`/healthcheck`) whose sole purpose is to signal that the server is running.
* No path parameters, multiple routes, or middleware chains are required for this use case.
* Using a third-party router would add an external dependency and conceptual overhead without tangible benefit.

**Alternatives Considered:**

* go-chi/chi (lightweight router with path parameters, middleware, sub-routers)
  → Rejected as overcomplicating for a single endpoint that only signals the server is running; net/http is sufficient

**Implications:**

* No router dependency in Stage 1
* If the HTTP surface grows in later stages (e.g. many endpoints or middleware needs), revisiting Chi or a similar router remains an option

---

## Final Note

These decisions reflect **intentional trade-offs**, not missing knowledge.

MyGoMetrics is designed to be:

* simple
* predictable
* observable
* production-appropriate

Rather than exhaustive or over-engineered.
