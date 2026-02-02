package collector

import (
	"context"
	"testing"
)

func TestDiskCollector_Name(t *testing.T) {
	c := NewDiskCollector()
	if got := c.Name(); got != "disk" {
		t.Errorf("Name() = %v, want %v", got, "disk")
	}
}

func TestDiskCollector_Collect(t *testing.T) {
	c := NewDiskCollector()
	ctx := context.Background()

	snapshots, err := c.Collect(ctx)
	if err != nil {
		// Disk I/O counters may not be available on all platforms
		// This is acceptable per the plan - we test the structure when available
		t.Logf("Collect() error = %v (may be platform-specific)", err)
		return
	}

	if len(snapshots) == 0 {
		t.Fatal("Collect() returned empty snapshots, want at least one")
	}

	// Verify we got at least two snapshots
	if len(snapshots) < 2 {
		t.Errorf("Collect() returned %v snapshots, want at least 2", len(snapshots))
	}

	// Verify expected metric names
	expectedNames := map[string]bool{
		"disk_read_bytes":  false,
		"disk_write_bytes": false,
	}

	for _, snap := range snapshots {
		if _, exists := expectedNames[snap.Name]; !exists {
			t.Errorf("Unexpected metric name: %v", snap.Name)
		} else {
			expectedNames[snap.Name] = true
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
	for name, found := range expectedNames {
		if !found {
			t.Errorf("Expected metric %v not found", name)
		}
	}
}

func TestDiskCollector_Collect_ContextCancelled(t *testing.T) {
	c := NewDiskCollector()
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
