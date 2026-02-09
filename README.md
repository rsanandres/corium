# Corium

Kubernetes operator built with Go and Kubebuilder, managing custom CRDs for automated stats collection, threshold-based alerting, and monitoring — deployed alongside a Next.js dashboard.

## Architecture

```
┌─────────────────────────────────────────┐
│           Kubernetes Cluster            │
│                                         │
│  ┌───────────────┐  ┌───────────────┐  │
│  │   Operator     │  │   Dashboard   │  │
│  │  (Go/K8s)     │  │   (Next.js)   │  │
│  └───────┬───────┘  └───────────────┘  │
│          │                              │
│  ┌───────▼──────────────────────────┐  │
│  │       Custom Resources           │  │
│  │  JAXStatsConfig                  │  │
│  │  JAXStatsCollector               │  │
│  │  JAXStatsAlert                   │  │
│  └──────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

## Custom Resource Definitions

API group: `stats.corium.io/v1alpha1`

### JAXStatsConfig

Global configuration for stats collection — enable/disable metrics, set collection intervals, configure storage backends (Prometheus, Elasticsearch).

### JAXStatsCollector

Collects stats from targeted pods via label selectors. Metrics include memory usage, GPU utilization, training metrics, and model performance. Supports cron-based collection schedules.

### JAXStatsAlert

Threshold-based alerting rules linked to collectors. Configurable severity levels, cooldown periods, and notification channels (Slack, email, webhook) with templated messages.

## Tech Stack

| Component | Technology |
|-----------|------------|
| Operator | Go 1.24, Kubebuilder v4.6, controller-runtime |
| Dashboard | Next.js 14, React 18, Tailwind CSS, Framer Motion |
| CRDs | `stats.corium.io/v1alpha1` |
| CI | GitHub Actions (lint, unit tests, e2e with KinD) |
| Deployment | Kustomize, Prometheus ServiceMonitor |

## Project Structure

```
corium/
├── operator/                # Go Kubernetes operator
│   ├── api/v1alpha1/        # CRD type definitions
│   ├── internal/controller/ # Reconciliation controllers
│   ├── config/              # RBAC, CRDs, Kustomize overlays
│   └── Makefile
├── src/                     # Next.js dashboard
├── k8s/                     # Deployment manifests
├── jaxstats/                # Stats app (submodule)
└── .github/workflows/       # CI pipelines
```

## Quick Start

```bash
# Install CRDs
cd operator && make install

# Run the operator locally
make run

# Deploy to cluster
make deploy
```

## License

MIT
