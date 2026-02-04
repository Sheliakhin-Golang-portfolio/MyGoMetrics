# Observability Guide

This document describes how to set up end-to-end observability for MyGoMetrics using Prometheus and Grafana.

---

## Overview

MyGoMetrics exposes metrics on the `/metrics` endpoint in Prometheus format. To build a complete observability stack:

1. **Prometheus** scrapes metrics from MyGoMetrics
2. **Prometheus** evaluates alert rules and sends notifications (via Alertmanager, if configured)
3. **Grafana** visualizes metrics using the provided dashboard

---

## Prerequisites

- MyGoMetrics running and accessible (see [RUNNING.md](RUNNING.md))
- Prometheus installed and configured
- Grafana installed and configured
- Prometheus datasource configured in Grafana

---

## Step 1: Configure Prometheus Scraping

Add MyGoMetrics to your Prometheus scrape configuration (`prometheus.yml`):

```yaml
scrape_configs:
  - job_name: 'mygometrics'
    static_configs:
      - targets: ['localhost:9000']  # Adjust to your MyGoMetrics address
    scrape_interval: 30s
    metrics_path: '/metrics'
```

If MyGoMetrics is running with custom `host` and `env` labels (via configuration), Prometheus will automatically collect these labels with the metrics.

**Verification:**

1. Start Prometheus
2. Navigate to `http://localhost:9090/targets`
3. Verify that the `mygometrics` target is `UP`
4. Query `up{job="mygometrics"}` in Prometheus UI — should return `1`

---

## Step 2: Configure Prometheus Alert Rules

### Loading Alert Rules

Add the alert rules file to your Prometheus configuration (`prometheus.yml`):

```yaml
rule_files:
  - '/path/to/prometheus/alerts.yml'  # Path to MyGoMetrics alerts.yml
```

If running Prometheus in Docker or Kubernetes, mount the `prometheus/alerts.yml` file and reference it appropriately.

### Alert Rules Included

The `prometheus/alerts.yml` file includes the following example alert rules:

- **MyGoMetricsExporterDown**: Fires when the exporter is unreachable for 2+ minutes
- **MyGoMetricsHighCPU**: Fires when CPU usage exceeds 90% for 5+ minutes
- **MyGoMetricsHighMemoryUsage**: Fires when memory usage exceeds 85% for 5+ minutes
- **MyGoMetricsHighGoroutineCount**: Fires when goroutine count exceeds 1000 for 5+ minutes

**Customization:**

- Adjust thresholds in `prometheus/alerts.yml` to match your environment
- Modify `for` durations to change alert sensitivity
- Update `job` label matching in the exporter down alert if your scrape job name differs

**Verification:**

1. Reload Prometheus configuration (or restart Prometheus)
2. Navigate to `http://localhost:9090/alerts`
3. Verify that alert rules appear and are evaluated

**Alertmanager Integration:**

To send alerts to Alertmanager, configure Prometheus with:

```yaml
alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - 'alertmanager:9093'  # Adjust to your Alertmanager address
```

---

## Step 3: Import Grafana Dashboard

### Import Steps

1. **Open Grafana** and navigate to Dashboards
2. Click **"New"** → **"Import"**
3. **Upload** the `grafana/dashboard.json` file (or paste its contents)
4. **Select** your Prometheus datasource from the dropdown
5. Click **"Import"**

The dashboard will be imported with the UID `mygometrics-overview` and title **"MyGoMetrics - Host & Runtime Overview"**.

### Dashboard Features

The dashboard includes:

- **CPU Usage**: Time series panel showing CPU usage percentage (0-100%)
- **Memory Usage %**: Gauge panel showing memory utilization percentage
- **Memory Usage**: Time series panel showing used and total memory in bytes
- **Disk I/O**: Time series panel showing disk read/write rates (bytes/second)
- **Go Runtime: Goroutines**: Time series panel showing current goroutine count
- **Go Runtime: Heap Allocation**: Time series panel showing heap memory allocation
- **Go Runtime: GC Cycles Rate**: Time series panel showing garbage collection cycle rate

### Dashboard Variables

The dashboard includes two template variables for filtering:

- **Host**: Filter metrics by `host` label (defaults to "All")
- **Environment**: Filter metrics by `env` label (defaults to "All")

These variables are populated automatically from available label values in Prometheus.

### Customization

- **Time Range**: Adjust the default time range (currently 6 hours) in dashboard settings
- **Refresh Interval**: Modify the auto-refresh interval (currently 30 seconds)
- **Panels**: Add, remove, or modify panels as needed
- **Thresholds**: Adjust panel thresholds and colors in panel settings

**Reference:** [Grafana Dashboard Import Documentation](https://grafana.com/docs/grafana/latest/dashboards/export-import/)

---

## End-to-End Observability Flow

```
┌─────────────┐
│ MyGoMetrics │
│  :9000     │
└──────┬──────┘
       │ HTTP GET /metrics
       │ (every 30s)
       ▼
┌─────────────┐
│ Prometheus  │
│  :9090      │
└──────┬──────┘
       │
       ├───► Scrapes metrics
       │     Stores time series
       │
       ├───► Evaluates alert rules
       │     (prometheus/alerts.yml)
       │     └───► Sends to Alertmanager
       │           (if configured)
       │
       └───► Serves metrics to Grafana
             (via PromQL queries)
                    │
                    ▼
            ┌─────────────┐
            │   Grafana   │
            │   :3000     │
            └─────────────┘
                  │
                  └───► Displays dashboard
                        (grafana/dashboard.json)
```

### Data Flow

1. **Collection**: MyGoMetrics collects system and runtime metrics periodically (configurable interval)
2. **Exposition**: Metrics are exposed on `/metrics` endpoint in Prometheus format
3. **Scraping**: Prometheus scrapes `/metrics` at configured intervals
4. **Storage**: Prometheus stores time series data
5. **Alerting**: Prometheus evaluates alert rules and triggers alerts when conditions are met
6. **Visualization**: Grafana queries Prometheus via PromQL and displays metrics in dashboards

---

## Troubleshooting

### Prometheus Not Scraping

- Verify MyGoMetrics is running: `curl http://localhost:9000/healthcheck`
- Check Prometheus targets: `http://localhost:9090/targets`
- Verify network connectivity between Prometheus and MyGoMetrics
- Check Prometheus logs for scrape errors

### Grafana Shows "No Data"

- Verify Prometheus datasource is configured correctly in Grafana
- Test datasource connection in Grafana datasource settings
- Verify metrics exist in Prometheus: query `mygometrics_cpu_usage_percent` in Prometheus UI
- Check dashboard time range — ensure it overlaps with available data
- Verify label filters in dashboard variables match your metric labels

### Alerts Not Firing

- Verify alert rules are loaded: check `http://localhost:9090/alerts`
- Ensure alert conditions are met (check current metric values)
- Verify `for` duration has elapsed
- Check Prometheus logs for rule evaluation errors
- If using Alertmanager, verify it's configured and reachable

---

## Additional Resources

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/grafana/latest/)
- [PromQL Guide](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [MyGoMetrics Configuration](CONFIGURATION.md)
- [MyGoMetrics Running Guide](RUNNING.md)
