package collector

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v3/mem"
)

// MemoryCollector collects memory metrics using gopsutil.
// It exposes used and total memory in bytes.
type MemoryCollector struct{}

// NewMemoryCollector creates a new MemoryCollector instance.
func NewMemoryCollector() *MemoryCollector {
	return &MemoryCollector{}
}

// Name returns the name of this collector.
func (c *MemoryCollector) Name() string {
	return "memory"
}

// Collect gathers memory metrics and returns them as Snapshots.
// Metrics collected:
//   - memory_used_bytes: amount of memory currently used
//   - memory_total_bytes: total amount of memory available
//
// If collection fails, returns an error with partial results (empty slice).
func (c *MemoryCollector) Collect(ctx context.Context) ([]Snapshot, error) {
	// Check if context is cancelled
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("memory: %w", err)
	}

	snapshots := []Snapshot{
		{
			Name:   "memory_used_bytes",
			Value:  float64(vmStat.Used),
			Labels: nil,
		},
		{
			Name:   "memory_total_bytes",
			Value:  float64(vmStat.Total),
			Labels: nil,
		},
	}

	return snapshots, nil
}
