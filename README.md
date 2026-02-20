# Corium

Production-style Kubernetes operator demonstrating real-world controller patterns — automated pod discovery, metrics collection, threshold-based alerting, and full observability with Prometheus + Grafana. Deployed alongside a Next.js dashboard in a multi-namespace cluster with network isolation.

## Architecture

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                            Kubernetes Cluster                                │
│                                                                              │
│  ┌──────────────── operator-system ────────────────────────────────────────┐ │
│  │                                                                          │ │
│  │  ┌──────────────── Operator Controller Manager ───────────────────────┐ │ │
│  │  │                                                                     │ │ │
│  │  │  ┌──────────────┐  ┌─────────────────┐  ┌────────────────────┐   │ │ │
│  │  │  │    Config     │  │    Collector     │  │      Alert         │   │ │ │
│  │  │  │  Controller   │  │   Controller     │  │    Controller      │   │ │ │
│  │  │  └──────┬────────┘  └────────┬─────────┘  └─────────┬─────────┘   │ │ │
│  │  │         │ reconciles         │ reconciles            │ reconciles  │ │ │
│  │  │         ▼                    ▼                       ▼             │ │ │
│  │  │  ┌──────────────┐  ┌─────────────────┐  ┌────────────────────┐   │ │ │
│  │  │  │ CoriumMonitor│◄─│ CoriumMonitor    │◄─│  CoriumMonitor     │   │ │ │
│  │  │  │  Config(CRD) │  │  Collector(CRD)  │  │   Alert(CRD)      │   │ │ │
│  │  │  │              │  │                   │  │                    │   │ │ │
│  │  │  │ • interval   │  │ • targetNamespace │  │ • rules[]         │   │ │ │
│  │  │  │ • metrics[]  │  │ • selector        │  │   metric/op/thold │   │ │ │
│  │  │  │ • storage    │  │ • configRef ──────┘  │ • collectorRef ───┘   │ │ │
│  │  │  └──────────────┘  └────────┬──────────┘  │ • cooldown           │ │ │
│  │  │                             │              └─────────┬────────────┘ │ │
│  │  │                             │ writes                 │ evaluates    │ │ │
│  │  │                   ┌─────────▼──────────┐   ┌────────▼──────────┐  │ │ │
│  │  │                   │    ConfigMap        │──▶│   K8s Events      │  │ │ │
│  │  │                   │  {name}-metrics     │   │ AlertFiring       │  │ │ │
│  │  │                   │  (metrics.json)     │   │ AlertResolved     │  │ │ │
│  │  │                   └────────────────────┘   └───────────────────┘  │ │ │
│  │  │                                                                     │ │ │
│  │  │  Prometheus Metrics (:8443)                                        │ │ │
│  │  │  • corium_discovered_pods  • corium_active_alerts                  │ │ │
│  │  │  • corium_reconcile_errors_total                                   │ │ │
│  │  └─────────────────────────────────┬───────────────────────────────────┘ │ │
│  └────────────────────────────────────┼────────────────────────────────────┘ │
│                                       │ scrapes                              │
│  ┌──────────────── corium-monitoring ─┼────────────────────────────────────┐ │
│  │                                    │                                     │ │
│  │  ┌───────────────────┐  ┌─────────▼──────────┐                        │ │
│  │  │      Grafana       │◄─│    Prometheus       │                        │ │
│  │  │  (18-panel dash)   │  │                      │                        │ │
│  │  │ • Pod Discovery    │  │ • ServiceMonitor     │                        │ │
│  │  │ • Alert History    │  │ • 60s scrape         │                        │ │
│  │  │ • Reconcile Perf   │  └──────────────────────┘                        │ │
│  │  │ • Queue Depth      │                                                   │ │
│  │  └───────────────────┘                                                   │ │
│  └──────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  ┌──────────────── corium-workloads ────────────────────────────────────────┐ │
│  │                                                                           │ │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐         ┌──────────────────┐    │ │
│  │  │  Pod 1   │  │  Pod 2   │  │  Pod 3   │  ...    │  Next.js Dash   │    │ │
│  │  │app:corium│  │app:corium│  │app:corium│         │  (3 replicas)   │    │ │
│  │  └─────────┘  └─────────┘  └─────────┘         └──────────────────┘    │ │
│  │       ▲              ▲             ▲                                      │ │
│  │       └──────────────┼─────────────┘                                     │ │
│  │            Collector discovers pods via label selector                    │ │
│  └───────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  NetworkPolicies: default-deny per namespace + explicit allows               │
│  • monitoring → operator :8443  • grafana → prometheus :9090                │
│  • intra-workloads allowed                                                  │
└──────────────────────────────────────────────────────────────────────────────┘
```

## How It Works

Three Custom Resource Definitions form a dependency chain, each managed by its own controller:

1. **CoriumMonitorConfig** defines what to collect (metrics list, interval, storage backend). Owns a finalizer that cascading-deletes dependent Collectors when removed.
2. **CoriumMonitorCollector** discovers pods via label selectors in a target namespace, collects their metrics, and persists results to an owned ConfigMap. Owner references ensure automatic cleanup.
3. **CoriumMonitorAlert** evaluates threshold rules against the Collector's metrics and emits Kubernetes Events (`AlertFiring` / `AlertResolved`) with configurable cooldown periods.

The operator exposes custom Prometheus metrics (`corium_discovered_pods`, `corium_active_alerts`, `corium_reconcile_errors_total`) scraped via ServiceMonitor. A pre-built 18-panel Grafana dashboard auto-provisions via ConfigMap sidecar.

## Kubernetes Patterns Demonstrated

| Pattern | Where |
|---------|-------|
| **Finalizers** | Config controller prevents orphaned Collectors on deletion |
| **Status Conditions** | `metav1.Condition` with `ObservedGeneration` on all CRDs |
| **Owner References** | Collector-owned ConfigMaps auto-deleted with parent CR |
| **EventRecorder** | Alert controller emits typed K8s Events for firing/resolved |
| **Cross-resource reconciliation** | Alert reads Collector's ConfigMap, Collector reads Config |
| **Kubebuilder validation markers** | Enum, MinLength, MinItems, Min/Max constraints on CRD fields |
| **Printer columns** | `kubectl get` shows Enabled, Status, Pods, Firing counts inline |
| **NetworkPolicies** | Default-deny per namespace with explicit ingress/egress rules |
| **Prometheus custom metrics** | 3 operator-level gauges + counter via controller-runtime |
| **Grafana dashboard** | 18-panel JSON auto-provisioned via ConfigMap sidecar |
| **ServiceMonitor** | Prometheus auto-discovers operator scrape target |
| **Namespace isolation** | 3 namespaces with distinct security boundaries |

## Custom Resource Definitions

API group: `monitor.corium.io/v1alpha1`

| CRD | Short Name | Purpose |
|-----|------------|---------|
| **CoriumMonitorConfig** | `cmc` | Global collection settings (interval, metrics, storage) |
| **CoriumMonitorCollector** | `cmcol` | Pod discovery + metrics persistence to ConfigMap |
| **CoriumMonitorAlert** | `cma` | Threshold alerting with K8s Events + cooldown |

## Tech Stack

| Component | Technology |
|-----------|------------|
| Operator | Go 1.24, Kubebuilder v4.6, controller-runtime v0.21 |
| Dashboard | Next.js 14, React 18, Tailwind CSS, Framer Motion |
| CRDs | `monitor.corium.io/v1alpha1` (3 resources) |
| Observability | Prometheus + Grafana (kube-prometheus-stack), 18-panel dashboard |
| Testing | Ginkgo v2 + Gomega, envtest (in-memory K8s API server) |
| CI | GitHub Actions (lint, type-check, unit tests, operator tests) |
| Deployment | Kustomize, KinD, NetworkPolicies, namespace isolation |

## Project Structure

```
corium/
├── operator/                    # Go Kubernetes operator
│   ├── api/v1alpha1/            # CRD type definitions with validation markers
│   ├── internal/controller/     # 3 reconciliation controllers + pure functions
│   │   ├── metrics.go           # Prometheus custom metrics registration
│   │   ├── metrics_collector.go # Pod metrics collection (pure, testable)
│   │   └── alert_evaluator.go   # Alert rule evaluation (pure, testable)
│   ├── config/                  # RBAC, CRDs, Kustomize, ServiceMonitor, NetworkPolicy
│   └── Makefile
├── app/                         # Next.js dashboard (App Router)
├── k8s/                         # Deployment manifests
│   ├── namespaces.yaml          # 3 isolated namespaces
│   ├── network-policies.yaml    # Default-deny + explicit allow rules
│   └── monitoring/              # Prometheus values + Grafana dashboard (18 panels)
├── demo.sh                      # One-command KinD demo
└── .github/workflows/           # CI pipelines
```

## Quick Start

```bash
# Full demo (requires Docker + KinD)
./demo.sh

# Or manually:
cd operator
make install    # Install CRDs
make run        # Run operator locally
kubectl apply -f config/samples/
kubectl get cmc,cmcol,cma
```

See [operator/README.md](operator/README.md) for detailed architecture diagrams and CRD reference.

## License

MIT
