package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

// metricType represents whether a metric is a Gauge or Counter.
type metricType int

const (
	metricTypeGauge metricType = iota
	metricTypeCounter
)

// metricDescriptor describes a Prometheus metric: its name, type, and help text.
type metricDescriptor struct {
	name     string
	typ      metricType
	help     string
	snapshot string // snapshot name from collector
}

// metricDescriptors maps collector snapshot names to Prometheus metric descriptors.
var metricDescriptors = map[string]metricDescriptor{
	"cpu_usage_percent": {
		name:     "mygometrics_cpu_usage_percent",
		typ:      metricTypeGauge,
		help:     "CPU usage percentage (0-100)",
		snapshot: "cpu_usage_percent",
	},
	"memory_used_bytes": {
		name:     "mygometrics_memory_used_bytes",
		typ:      metricTypeGauge,
		help:     "Memory used in bytes",
		snapshot: "memory_used_bytes",
	},
	"memory_total_bytes": {
		name:     "mygometrics_memory_total_bytes",
		typ:      metricTypeGauge,
		help:     "Total memory available in bytes",
		snapshot: "memory_total_bytes",
	},
	"disk_read_bytes": {
		name:     "mygometrics_disk_read_bytes",
		typ:      metricTypeCounter,
		help:     "Total bytes read from disk (cumulative)",
		snapshot: "disk_read_bytes",
	},
	"disk_write_bytes": {
		name:     "mygometrics_disk_write_bytes",
		typ:      metricTypeCounter,
		help:     "Total bytes written to disk (cumulative)",
		snapshot: "disk_write_bytes",
	},
	"runtime_goroutines": {
		name:     "mygometrics_runtime_goroutines",
		typ:      metricTypeGauge,
		help:     "Current number of goroutines",
		snapshot: "runtime_goroutines",
	},
	"runtime_gc_cycles": {
		name:     "mygometrics_runtime_gc_cycles",
		typ:      metricTypeCounter,
		help:     "Total number of GC cycles since program start (cumulative)",
		snapshot: "runtime_gc_cycles",
	},
	"runtime_heap_alloc_bytes": {
		name:     "mygometrics_runtime_heap_alloc_bytes",
		typ:      metricTypeGauge,
		help:     "Bytes allocated on the heap",
		snapshot: "runtime_heap_alloc_bytes",
	},
}

// metricHolders holds all Prometheus metrics (gauges and counters) keyed by snapshot name.
type metricHolders struct {
	gauges   map[string]*prometheus.GaugeVec
	counters map[string]*prometheus.CounterVec
	// For counters, we need to track previous values to compute deltas
	prevValues map[string]float64
}

// createMetrics creates and registers all Prometheus metrics with the given registry.
// Labels (host, env) are applied as variable labels when updating metrics in Registry.Update.
func createMetrics(reg *prometheus.Registry) (*metricHolders, error) {
	holders := &metricHolders{
		gauges:     make(map[string]*prometheus.GaugeVec),
		counters:   make(map[string]*prometheus.CounterVec),
		prevValues: make(map[string]float64),
	}

	labelNames := []string{"host", "env"}

	for snapshotName, desc := range metricDescriptors {
		opts := prometheus.Opts{
			Name: desc.name,
			Help: desc.help,
		}

		switch desc.typ {
		case metricTypeGauge:
			gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts(opts), labelNames)
			if err := reg.Register(gauge); err != nil {
				return nil, err
			}
			holders.gauges[snapshotName] = gauge

		case metricTypeCounter:
			counter := prometheus.NewCounterVec(prometheus.CounterOpts(opts), labelNames)
			if err := reg.Register(counter); err != nil {
				return nil, err
			}
			holders.counters[snapshotName] = counter
		}
	}

	return holders, nil
}

// updateMetric updates the Prometheus metric corresponding to the given snapshot.
// For counters, it computes and adds the delta from the previous value.
func (h *metricHolders) updateMetric(snapshotName string, value float64, labelValues []string) {
	desc, exists := metricDescriptors[snapshotName]
	if !exists {
		// Unknown snapshot name, skip
		return
	}

	switch desc.typ {
	case metricTypeGauge:
		if gauge, ok := h.gauges[snapshotName]; ok {
			gauge.WithLabelValues(labelValues...).Set(value)
		}
	case metricTypeCounter:
		if counter, ok := h.counters[snapshotName]; ok {
			prevValue, hasPrev := h.prevValues[snapshotName]
			if !hasPrev {
				// First value: add the full value
				counter.WithLabelValues(labelValues...).Add(value)
			} else {
				// Subsequent values: add delta
				delta := value - prevValue
				if delta >= 0 {
					counter.WithLabelValues(labelValues...).Add(delta)
				}
				// If delta < 0, counter was reset (e.g., process restart), ignore
			}
			h.prevValues[snapshotName] = value
		}
	}
}
