package collector

import "context"

// Snapshot represents a single metric snapshot with its name, value, and optional labels.
// This is a Prometheus-agnostic representation that can be converted to Prometheus
// metrics in a later stage.
//
// Labels may be nil to indicate no labels. Callers should check `if s.Labels != nil`
// before ranging over the map.
type Snapshot struct {
	// The metric name (e.g., "cpu_usage_percent", "memory_used_bytes").
	Name string

	Value float64

	// Labels is an optional map of label key-value pairs.
	// If nil, the metric has no labels.
	Labels map[string]string
}

// Collector is the interface that all metric collectors must implement.
// Collectors gather metrics from various sources (system, runtime, etc.)
// and return them as Snapshot slices.
//
// Implementations should handle errors gracefully and return partial results
// rather than failing completely, as per the project's error handling strategy.
type Collector interface {
	// Name returns the name of this collector (e.g., "cpu", "memory", "disk", "runtime").
	Name() string

	// Collect gathers metrics and returns them as a slice of Snapshots.
	// The context can be used for cancellation and timeouts.
	// If an error occurs, implementations should return partial results along
	// with the error, allowing callers to process available metrics.
	Collect(ctx context.Context) ([]Snapshot, error)
}
