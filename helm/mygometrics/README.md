# MyGoMetrics Helm Chart

This Helm chart deploys MyGoMetrics, a host-level Prometheus exporter for system and Go runtime metrics, on Kubernetes.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- Prometheus Operator (optional, for ServiceMonitor support)

## Installation

### Install from local chart

```bash
helm install mygometrics ./helm/mygometrics
```

### Install with custom values

```bash
helm install mygometrics ./helm/mygometrics \
  --set image.repository=your-registry/mygometrics \
  --set image.tag=0.6.0 \
  --set config.env=production
```

### Install with values file

```bash
helm install mygometrics ./helm/mygometrics -f my-values.yaml
```

## Configuration

See `values.yaml` for all available configuration options. Key settings:

- **Image**: `image.repository`, `image.tag`, `image.pullPolicy`
- **Config**: All application settings under `config.*` (mapped to environment variables)
- **Service**: `service.type`, `service.port`
- **Resources**: `resources.requests`, `resources.limits`
- **ServiceMonitor**: `serviceMonitor.enabled`, `serviceMonitor.interval`, `serviceMonitor.labels`

For detailed configuration options, see the application's [CONFIGURATION.md](../../docs/CONFIGURATION.md).

## ServiceMonitor

To enable Prometheus Operator integration, set `serviceMonitor.enabled: true` and configure `serviceMonitor.labels` to match your Prometheus instance's selector:

```yaml
serviceMonitor:
  enabled: true
  interval: "30s"
  labels:
    release: prometheus
```

## Upgrading

```bash
helm upgrade mygometrics ./helm/mygometrics
```

## Uninstallation

```bash
helm uninstall mygometrics
```
