package collector

import (
	"context"
	"testing"
)

func TestRuntimeCollector_Name(t *testing.T) {
	c := NewRuntimeCollector()
	if got := c.Name(); got != "runtime" {
		t.Errorf("Name() = %v, want %v", got, "runtime")
	}
}

func TestRuntimeCollector_Collect(t *testing.T) {
	c := NewRuntimeCollector()
	ctx := context.Background()

	snapshots, err := c.Collect(ctx)
	if err != nil {
		t.Fatalf("Collect() error = %v, want nil", err)
	}

	if len(snapshots) == 0 {
		t.Fatal("Collect() returned empty snapshots, want at least one")
	}

	// Verify we got exactly three snapshots
	actualNumberOfSnapshots := 0

	goroutinesSnapIndex := -1

	for index, snap := range snapshots {
		switch snap.Name {
		case "runtime_goroutines":
			goroutinesSnapIndex = index
			actualNumberOfSnapshots++
		case "runtime_gc_cycles":
			actualNumberOfSnapshots++
		case "runtime_heap_alloc_bytes":
			actualNumberOfSnapshots++
		default:
			t.Errorf("Unexpected metric name: %v", snap.Name)
		}

		// Verify value is non-negative
		if snap.Value < 0 {
			t.Errorf("Metric %v has negative value: %v", snap.Name, snap.Value)
		}

		// Verify labels are nil (as per implementation)
		if snap.Labels != nil {
			t.Errorf("Metric %v has non-nil labels, want nil", snap.Name)
		}
	}

	// Verify all expected metrics are present
	if actualNumberOfSnapshots != 3 {
		t.Errorf("Expected 3 specific snapshots, got %v", actualNumberOfSnapshots)
	}

	// Verify goroutines count is at least 1 (the test itself runs in a goroutine)
	if snapshots[goroutinesSnapIndex].Value < 1 {
		t.Errorf("runtime_goroutines = %v, want >= 1", snapshots[goroutinesSnapIndex].Value)
	}
}

func TestRuntimeCollector_Collect_ContextCancelled(t *testing.T) {
	c := NewRuntimeCollector()
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
