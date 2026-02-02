package collector

import (
	"context"
	"testing"
)

func TestCPUCollector_Name(t *testing.T) {
	c := NewCPUCollector()
	if got := c.Name(); got != "cpu" {
		t.Errorf("Name() = %v, want %v", got, "cpu")
	}
}

func TestCPUCollector_Collect(t *testing.T) {
	c := NewCPUCollector()
	ctx := context.Background()

	snapshots, err := c.Collect(ctx)
	if err != nil {
		t.Fatalf("Collect() error = %v, want nil", err)
	}

	if len(snapshots) == 0 {
		t.Fatal("Collect() returned empty snapshots, want at least one")
	}

	// Verify we got exactly one snapshot
	if len(snapshots) != 1 {
		t.Errorf("Collect() returned %v snapshots, want 1", len(snapshots))
	}

	snap := snapshots[0]

	// Verify metric name
	if snap.Name != "cpu_usage_percent" {
		t.Errorf("Metric name = %v, want %v", snap.Name, "cpu_usage_percent")
	}

	// Verify value is in valid range (0-100 for percentage)
	if snap.Value < 0 || snap.Value > 100 {
		t.Errorf("CPU usage = %v, want value between 0 and 100", snap.Value)
	}

	// Verify labels are nil
	if snap.Labels != nil {
		t.Errorf("Labels = %v, want nil", snap.Labels)
	}
}

func TestCPUCollector_Collect_ContextCancelled(t *testing.T) {
	c := NewCPUCollector()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	snapshots, err := c.Collect(ctx)
	if err == nil {
		t.Error("Collect() error = nil, want context cancelled error")
	}
	if snapshots != nil {
		t.Errorf("Collect() returned snapshots on cancelled context: %v", snapshots)
	}
}
