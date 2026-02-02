package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

// CPUCollector collects CPU usage metrics using gopsutil.
// It exposes the total CPU usage percentage.
type CPUCollector struct{}

// NewCPUCollector creates a new CPUCollector instance.
func NewCPUCollector() *CPUCollector {
	return &CPUCollector{}
}

// Name returns the name of this collector.
func (c *CPUCollector) Name() string {
	return "cpu"
}

// Collect gathers CPU metrics and returns them as Snapshots.
// Metrics collected:
//   - cpu_usage_percent: total CPU usage percentage (0-100)
//
// Uses a 1-second interval to measure CPU usage. If collection fails,
// returns an error with partial results (empty slice).
func (c *CPUCollector) Collect(ctx context.Context) ([]Snapshot, error) {
	// Check if context is cancelled
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Use PercentWithContext to respect context cancellation
	// Interval of 1 second is reasonable for CPU usage measurement
	percentages, err := cpu.PercentWithContext(ctx, 1*time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("cpu: %w", err)
	}

	if len(percentages) == 0 {
		return nil, fmt.Errorf("cpu: no CPU data available")
	}

	// Use the first (and typically only) value for total CPU usage
	snapshots := []Snapshot{
		{
			Name:   "cpu_usage_percent",
			Value:  percentages[0],
			Labels: nil,
		},
	}

	return snapshots, nil
}
