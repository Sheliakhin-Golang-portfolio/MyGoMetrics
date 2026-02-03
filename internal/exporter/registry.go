package exporter

import (
	"context"
	"net/http"

	"github.com/Sheliakhin-Golang-portfolio/MyGoMetrics/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Registry bridges collector snapshots to Prometheus metrics.
// It iterates collectors periodically and updates Prometheus Gauges/Counters.
type Registry struct {
	collectors []collector.Collector
	reg        *prometheus.Registry
	metrics    *metricHolders
	labels     map[string]string
	logger     *zap.Logger
}

// NewRegistry creates a Registry with the given collectors, label values, and logger.
// Labels must include "host" and "env" (use empty string if not set).
// The Prometheus registry and all metrics are created and registered.
func NewRegistry(collectors []collector.Collector, host, env string, logger *zap.Logger) (*Registry, error) {
	labels := map[string]string{
		"host": host,
		"env":  env,
	}

	reg := prometheus.NewRegistry()
	metrics, err := createMetrics(reg)
	if err != nil {
		return nil, err
	}

	return &Registry{
		collectors: collectors,
		reg:        reg,
		metrics:    metrics,
		labels:     labels,
		logger:     logger,
	}, nil
}

// Update runs all collectors and updates Prometheus metrics.
// Per-collector errors are logged; partial results are still written.
// This method should be called periodically (e.g., every CollectInterval).
func (r *Registry) Update(ctx context.Context) {
	labelValues := []string{r.labels["host"], r.labels["env"]}

	for _, c := range r.collectors {
		snapshots, err := c.Collect(ctx)
		if err != nil {
			r.logger.Warn("collector failed", zap.String("collector", c.Name()), zap.Error(err))
			// Continue with other collectors per error handling strategy
			continue
		}

		// Update Prometheus metrics for each snapshot
		for _, snapshot := range snapshots {
			r.metrics.updateMetric(snapshot.Name, snapshot.Value, labelValues)
		}
	}
}

// Handler returns an http.Handler for the /metrics endpoint.
// The handler serves the Prometheus registry using promhttp.HandlerFor.
func (r *Registry) Handler() http.Handler {
	return promhttp.HandlerFor(r.reg, promhttp.HandlerOpts{})
}
