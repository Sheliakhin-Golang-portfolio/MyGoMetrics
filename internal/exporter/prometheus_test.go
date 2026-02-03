package exporter

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Sheliakhin-Golang-portfolio/MyGoMetrics/internal/collector"
	"go.uber.org/zap"
)

// mockCollectorSuccess is a mock collector that always succeeds
type mockCollectorSuccess struct {
	name      string
	snapshots []collector.Snapshot
}

func (m *mockCollectorSuccess) Name() string {
	return m.name
}

func (m *mockCollectorSuccess) Collect(ctx context.Context) ([]collector.Snapshot, error) {
	return m.snapshots, nil
}

// mockCollectorFailure is a mock collector that always fails
type mockCollectorFailure struct {
	name string
	err  error
}

func (m *mockCollectorFailure) Name() string {
	return m.name
}

func (m *mockCollectorFailure) Collect(ctx context.Context) ([]collector.Snapshot, error) {
	return nil, m.err
}

func TestRegistry_Update_OneFailingCollectorDoesNotBreakMetrics(t *testing.T) {
	tests := []struct {
		name           string
		collectors     []collector.Collector
		expectedStatus int
		shouldContain  string // metric name that should be present
	}{
		{
			name: "one failing collector with one successful",
			collectors: []collector.Collector{
				&mockCollectorSuccess{
					name: "success",
					snapshots: []collector.Snapshot{
						{Name: "cpu_usage_percent", Value: 50.0},
					},
				},
				&mockCollectorFailure{
					name: "failure",
					err:  errors.New("mock error"),
				},
			},
			expectedStatus: http.StatusOK,
			shouldContain:  "mygometrics_cpu_usage_percent",
		},
		{
			name: "all collectors succeed",
			collectors: []collector.Collector{
				&mockCollectorSuccess{
					name: "cpu",
					snapshots: []collector.Snapshot{
						{Name: "cpu_usage_percent", Value: 75.0},
					},
				},
				&mockCollectorSuccess{
					name: "memory",
					snapshots: []collector.Snapshot{
						{Name: "memory_used_bytes", Value: 1024.0},
					},
				},
			},
			expectedStatus: http.StatusOK,
			shouldContain:  "mygometrics_cpu_usage_percent",
		},
		{
			name: "all collectors fail",
			collectors: []collector.Collector{
				&mockCollectorFailure{
					name: "failure1",
					err:  errors.New("error 1"),
				},
				&mockCollectorFailure{
					name: "failure2",
					err:  errors.New("error 2"),
				},
			},
			expectedStatus: http.StatusOK, // Should still return 200
			shouldContain:  "",            // No custom metrics expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create registry with no-op logger for tests
			registry, err := NewRegistry(tt.collectors, "test-host", "test-env", zap.NewNop())
			if err != nil {
				t.Fatalf("Failed to create registry: %v", err)
			}

			// Update registry (triggers collection)
			ctx := context.Background()
			registry.Update(ctx)

			// Create HTTP request
			req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
			rec := httptest.NewRecorder()

			// Serve metrics endpoint
			handler := registry.Handler()
			handler.ServeHTTP(rec, req)

			// Assert status code
			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			// Assert response body contains expected metric (if specified)
			if tt.shouldContain != "" {
				body := rec.Body.String()
				if !strings.Contains(body, tt.shouldContain) {
					t.Errorf("Expected response body to contain %q, got: %s", tt.shouldContain, body)
				}
			}
		})
	}
}

func TestRegistry_Update_PartialResultsOnFailure(t *testing.T) {
	// Test that when one collector fails, metrics from successful collectors are still present
	successCollector := &mockCollectorSuccess{
		name: "success",
		snapshots: []collector.Snapshot{
			{Name: "memory_used_bytes", Value: 2048.0},
			{Name: "memory_total_bytes", Value: 4096.0},
		},
	}

	failureCollector := &mockCollectorFailure{
		name: "failure",
		err:  errors.New("collection failed"),
	}

	registry, err := NewRegistry([]collector.Collector{successCollector, failureCollector}, "test-host", "test-env", zap.NewNop())
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	ctx := context.Background()
	registry.Update(ctx)

	// Verify metrics endpoint still works
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	handler := registry.Handler()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	// Should contain metrics from successful collector
	if !strings.Contains(body, "mygometrics_memory_used_bytes") {
		t.Errorf("Expected response to contain memory_used_bytes metric, got: %s", body)
	}
	if !strings.Contains(body, "mygometrics_memory_total_bytes") {
		t.Errorf("Expected response to contain memory_total_bytes metric, got: %s", body)
	}
}
