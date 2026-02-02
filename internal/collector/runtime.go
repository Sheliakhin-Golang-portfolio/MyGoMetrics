package collector

import (
	"context"
	"runtime"
)

// RuntimeCollector collects Go runtime metrics using the standard library.
// It exposes goroutine count, GC cycles, and heap allocation metrics.
type RuntimeCollector struct{}

// NewRuntimeCollector creates a new RuntimeCollector instance.
func NewRuntimeCollector() *RuntimeCollector {
	return &RuntimeCollector{}
}

// Name returns the name of this collector.
func (c *RuntimeCollector) Name() string {
	return "runtime"
}

// Collect gathers Go runtime metrics and returns them as Snapshots.
// Metrics collected:
//   - runtime_goroutines: current number of goroutines
//   - runtime_gc_cycles: total number of GC cycles since program start
//   - runtime_heap_alloc_bytes: bytes allocated on the heap
//
// This implementation uses runtime.ReadMemStats for compatibility.
// Future versions may use runtime/metrics for consistency with Go 1.25+.
func (c *RuntimeCollector) Collect(ctx context.Context) ([]Snapshot, error) {
	// Check if context is cancelled
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	snapshots := []Snapshot{
		{
			Name:   "runtime_goroutines",
			Value:  float64(runtime.NumGoroutine()),
			Labels: nil,
		},
		{
			Name:   "runtime_gc_cycles",
			Value:  float64(memStats.NumGC),
			Labels: nil,
		},
		{
			Name:   "runtime_heap_alloc_bytes",
			Value:  float64(memStats.HeapAlloc),
			Labels: nil,
		},
	}

	return snapshots, nil
}
