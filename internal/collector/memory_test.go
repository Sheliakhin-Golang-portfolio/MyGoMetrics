package collector

import (
	"context"
	"testing"
)

func TestMemoryCollector_Name(t *testing.T) {
	c := NewMemoryCollector()
	if got := c.Name(); got != "memory" {
		t.Errorf("Name() = %v, want %v", got, "memory")
	}
}

func TestMemoryCollector_Collect(t *testing.T) {
	c := NewMemoryCollector()
	ctx := context.Background()

	snapshots, err := c.Collect(ctx)
	if err != nil {
		t.Fatalf("Collect() error = %v, want nil", err)
	}

	if len(snapshots) == 0 {
		t.Fatal("Collect() returned empty snapshots, want at least one")
	}

	// Verify we got at least two snapshots
	if len(snapshots) < 2 {
		t.Errorf("Collect() returned %v snapshots, want at least 2", len(snapshots))
	}

	// Verify we got exactly two expected snapshots
	actualNumberOfSnapshots := 0
	usedSnapIndex := -1
	totalSnapIndex := -1

	for index, snap := range snapshots {
		switch snap.Name {
		case "memory_used_bytes":
			usedSnapIndex = index
			actualNumberOfSnapshots++
		case "memory_total_bytes":
			totalSnapIndex = index
			actualNumberOfSnapshots++
		default:
			t.Errorf("Unexpected metric name: %v", snap.Name)
		}

		// Verify value is non-negative
		if snap.Value < 0 {
			t.Errorf("Metric %v has negative value: %v", snap.Name, snap.Value)
		}

		// Verify labels are nil
		if snap.Labels != nil {
			t.Errorf("Metric %v has non-nil labels, want nil", snap.Name)
		}
	}

	// Verify all expected metrics are present
	if actualNumberOfSnapshots != 2 {
		t.Errorf("Expected 2 specific snapshots, got %v", actualNumberOfSnapshots)
	}

	// Verify used <= total
	if snapshots[usedSnapIndex].Value > snapshots[totalSnapIndex].Value {
		t.Errorf("memory_used_bytes (%v) > memory_total_bytes (%v)", snapshots[usedSnapIndex].Value, snapshots[totalSnapIndex].Value)
	}
}

func TestMemoryCollector_Collect_ContextCancelled(t *testing.T) {
	c := NewMemoryCollector()
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
