package collector

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v3/disk"
)

// DiskCollector collects disk I/O metrics using gopsutil.
// It exposes aggregate read and write bytes across all devices.
type DiskCollector struct{}

// NewDiskCollector creates a new DiskCollector instance.
func NewDiskCollector() *DiskCollector {
	return &DiskCollector{}
}

// Name returns the name of this collector.
func (c *DiskCollector) Name() string {
	return "disk"
}

// Collect gathers disk I/O metrics and returns them as Snapshots.
// Metrics collected:
//   - disk_read_bytes: total bytes read across all devices
//   - disk_write_bytes: total bytes written across all devices
//
// This implementation aggregates I/O counters across all devices for simplicity.
// If collection fails (e.g., on unsupported platforms), returns an error
// with partial results (empty slice).
func (c *DiskCollector) Collect(ctx context.Context) ([]Snapshot, error) {
	// Check if context is cancelled
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Get I/O counters for all devices
	ioCounters, err := disk.IOCountersWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("disk: %w", err)
	}

	// Aggregate read and write bytes across all devices
	var totalReadBytes, totalWriteBytes uint64
	for _, counter := range ioCounters {
		totalReadBytes += counter.ReadBytes
		totalWriteBytes += counter.WriteBytes
	}

	snapshots := []Snapshot{
		{
			Name:   "disk_read_bytes",
			Value:  float64(totalReadBytes),
			Labels: nil,
		},
		{
			Name:   "disk_write_bytes",
			Value:  float64(totalWriteBytes),
			Labels: nil,
		},
	}

	return snapshots, nil
}
