# Corium Operator

Kubernetes operator for automated pod discovery, metrics collection, threshold-based alerting, and monitoring. Built with Go 1.24, Kubebuilder v4.6, and controller-runtime v0.21. Includes full observability with Prometheus custom metrics and an 18-panel Grafana dashboard.

## Architecture

### System Overview

```mermaid
flowchart TB
    subgraph K8s["Kubernetes Cluster"]
        subgraph ONS["operator-system"]
            OP[Operator Controller Manager]
        end

        subgraph WNS["corium-workloads"]
            D[Next.js Dashboard]
            WP1[Workload Pod 1]
            WP2[Workload Pod 2]
            WP3[Workload Pod N]
        end

        subgraph MNS["corium-monitoring"]
            PROM[Prometheus]
            GRAF["Grafana (18 panels)"]
        end

        subgraph CRDs["Custom Resources"]
            CFG[CoriumMonitorConfig]
            COL[CoriumMonitorCollector]
            ALT[CoriumMonitorAlert]
        end

        CM[ConfigMap\nmetrics.json]
        EV[K8s Events]
    end

    OP -->|reconciles| CFG
    OP -->|reconciles| COL
    OP -->|reconciles| ALT
    COL -->|discovers| WP1 & WP2 & WP3
    COL -->|writes| CM
    ALT -->|reads| CM
    ALT -->|emits| EV
    PROM -->|scrapes :8443| OP
    GRAF -->|queries| PROM
```

### CRD Dependency Chain

The three CRDs form a directed dependency graph — each resource references the one before it:

```mermaid
graph LR
    CFG["CoriumMonitorConfig\n(global settings)"]
    COL["CoriumMonitorCollector\n(pod discovery)"]
    ALT["CoriumMonitorAlert\n(threshold rules)"]
    CM["ConfigMap\n(pod metrics)"]
    EV["K8s Events\n(alert notifications)"]

    CFG -->|referenced by| COL
    COL -->|referenced by| ALT
    COL -->|creates/owns| CM
    ALT -->|reads| CM
    ALT -->|emits| EV
```

### Collector Reconciliation Flow

Each reconciliation loop fetches the CR, validates its referenced Config, discovers pods, collects metrics, persists to a ConfigMap, and updates Prometheus gauges:

```mermaid
sequenceDiagram
    participant CR as Collector CR Event
    participant R as Reconciler
    participant API as K8s API
    participant CM as ConfigMap

    CR->>R: Reconcile triggered
    R->>API: Get CoriumMonitorCollector
    R->>API: Get referenced CoriumMonitorConfig
    alt Config not found or disabled
        R->>API: Update status (Error/Disabled)
        R-->>CR: Requeue after 30s/60s
    end
    R->>API: List Pods by label selector
    R->>R: CollectPodMetrics(pods)
    R->>R: BuildCollectedMetrics()
    R->>R: MarshalMetrics() to JSON
    R->>CM: Create or Update ConfigMap
    R->>API: Update Collector status
    R->>R: Update Prometheus gauge
    R-->>CR: Requeue after interval
```

### Alert Evaluation Flow

Alerts evaluate threshold rules against collected metrics, enforce cooldown periods, and emit Kubernetes Events:

```mermaid
flowchart TD
    START([Reconcile triggered]) --> FETCH[Fetch CoriumMonitorAlert]
    FETCH --> ENABLED{Alert enabled?}
    ENABLED -->|No| DISABLED[Set status: Disabled]
    DISABLED --> DONE([Requeue])

    ENABLED -->|Yes| GETCOL[Fetch referenced Collector]
    GETCOL --> COLOK{Collector found?}
    COLOK -->|No| ERROR1[Set status: Error]
    ERROR1 --> DONE

    COLOK -->|Yes| GETCM[Fetch metrics ConfigMap]
    GETCM --> CMOK{ConfigMap found?}
    CMOK -->|No| PENDING[Set status: Pending]
    PENDING --> DONE

    CMOK -->|Yes| PARSE[Parse metrics JSON]
    PARSE --> EVAL[EvaluateAlertRules]
    EVAL --> COOL{Cooldown elapsed?}
    COOL -->|Yes| EMIT[Emit K8s Events\nfor firing/resolved]
    COOL -->|No| SKIP[Skip event emission]
    EMIT --> STATUS
    SKIP --> STATUS[Update status\nFiring / OK]
    STATUS --> PROM[Update Prometheus gauge]
    PROM --> DONE
```

### Resource Lifecycle States

Every CRD follows a consistent state machine with well-defined transitions:

```mermaid
stateDiagram-v2
    [*] --> Pending: CR created
    Pending --> Active: Config valid + enabled
    Pending --> Disabled: Config disabled
    Pending --> Error: Config invalid / not found
    Active --> Error: Reconcile failure
    Active --> Disabled: Config disabled
    Error --> Active: Issue resolved
    Disabled --> Active: Config re-enabled
    Active --> [*]: CR deleted (finalizer cleanup)
```

## CRD Reference

### CoriumMonitorConfig

Global configuration that controls what metrics to collect and how often. Owns a finalizer (`monitor.corium.io/config-cleanup`) that cascading-deletes all dependent Collectors when the Config is removed.

| Field | Type | Description |
|-------|------|-------------|
| `spec.enabled` | bool | Enable/disable collection |
| `spec.collectionInterval` | int32 | Collection interval in seconds (10-3600) |
| `spec.metrics` | []string | Metrics to collect |
| `spec.storageConfig.type` | enum | Storage backend: `prometheus`, `configmap`, `elasticsearch` |

### CoriumMonitorCollector

Discovers pods via label selectors in a target namespace, collects their metrics, and persists results to an owned ConfigMap (`{name}-metrics`). Uses owner references so the ConfigMap is automatically garbage-collected when the Collector CR is deleted.

| Field | Type | Description |
|-------|------|-------------|
| `spec.targetNamespace` | string | Namespace to discover pods in |
| `spec.selector` | LabelSelector | Pod label selector |
| `spec.configRef` | string | Name of referenced CoriumMonitorConfig |
| `status.discoveredPods` | int32 | Number of discovered pods |
| `status.metricsConfigMap` | string | Name of the managed ConfigMap |

### CoriumMonitorAlert

Evaluates threshold rules against the Collector's metrics ConfigMap and emits typed Kubernetes Events (`AlertFiring` / `AlertResolved`). Supports configurable cooldown periods to prevent event flooding.

| Field | Type | Description |
|-------|------|-------------|
| `spec.enabled` | bool | Enable/disable alert evaluation |
| `spec.collectorRef` | string | Name of referenced CoriumMonitorCollector |
| `spec.cooldownPeriod` | string | Cooldown between alerts (e.g., "5m") |
| `spec.rules[].metric` | enum | `restart_count`, `not_ready_count`, `container_count`, `pod_count` |
| `spec.rules[].operator` | enum | `>`, `<`, `>=`, `<=`, `==` |
| `spec.rules[].threshold` | string | Numeric threshold value |
| `spec.rules[].severity` | enum | `critical`, `warning`, `info` |
| `status.firingAlertsCount` | int32 | Number of currently firing alerts |

## Kubernetes Patterns Demonstrated

| Pattern | Implementation |
|---------|---------------|
| **Finalizers** | Config controller prevents orphaned Collectors on deletion via `monitor.corium.io/config-cleanup` |
| **Status Conditions** | Standard `metav1.Condition` with `ObservedGeneration` tracking on all 3 CRDs |
| **Owner References** | Collector-owned ConfigMaps auto-deleted on CR removal (K8s garbage collection) |
| **EventRecorder** | Alert controller emits typed AlertFiring/AlertResolved events visible via `kubectl get events` |
| **Cross-resource reconciliation** | Alert reads Collector's ConfigMap; Collector reads Config — forming a dependency chain |
| **Kubebuilder validation markers** | Enum, MinLength, MinItems, Min/Max constraints enforced at admission time |
| **Printer columns** | `kubectl get` shows Enabled, Status, Pods, Firing counts inline without `-o yaml` |
| **NetworkPolicies** | Default-deny per namespace with explicit allow rules for cross-namespace communication |
| **Prometheus custom metrics** | `corium_discovered_pods`, `corium_active_alerts`, `corium_reconcile_errors_total` |
| **Grafana dashboard** | 18-panel dashboard auto-provisioned via ConfigMap sidecar (5 organized rows) |
| **ServiceMonitor** | Prometheus auto-discovers operator scrape target in `operator-system` namespace |
| **Pure function extraction** | Metrics collection and alert evaluation logic extracted into testable pure functions |

## Observability

### Prometheus Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `corium_discovered_pods` | Gauge | Pods currently discovered per Collector |
| `corium_active_alerts` | Gauge | Currently firing alerts per Alert resource |
| `corium_reconcile_errors_total` | Counter | Reconciliation errors per controller |

### Grafana Dashboard (18 panels, 5 rows)

| Row | Panels |
|-----|--------|
| **Overview** | Discovered Pods (stat), Active Alerts (stat), Reconciliation Rate (stat), Error Rate (stat) |
| **Pod Discovery & Collection** | Discovered Pods by Collector (bar gauge), Discovered Pods Over Time (timeseries) |
| **Alerting** | Active Alerts by Resource (gauge), Alert Firing History (timeseries) |
| **Reconciliation Performance** | Errors by Controller (timeseries), Success Rate % (timeseries), Duration p50/p99 (timeseries), Work Queue Depth (timeseries) |
| **Controller Runtime** | Reconciliations/s by Controller (bar chart), Work Queue Latency p99 (timeseries) |

## Testing

- **Framework:** Ginkgo v2 + Gomega (BDD-style)
- **Environment:** envtest bootstraps an in-memory K8s API server shared across all controller tests
- **Pattern:** Tests create CRs in `BeforeEach`, reconcile directly via `reconciler.Reconcile()`, assert status and side-effects, clean up in `AfterEach`
- **Event assertions:** `record.FakeRecorder` captures and verifies K8s events
- **Pure function tests:** `metrics_collector.go` and `alert_evaluator.go` tested independently

## Quick Start

```bash
# Run tests
make test

# Install CRDs
make install

# Run operator locally
make run

# Apply sample resources
kubectl apply -f config/samples/

# Check results
kubectl get cmc,cmcol,cma
kubectl get configmap -l app.kubernetes.io/managed-by=corium-operator
kubectl get events --field-selector reason=AlertFiring
```

## Full Demo

```bash
# One-command demo with KinD
cd .. && ./demo.sh
```

## License

MIT
