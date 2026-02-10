# Corium Operator

Kubernetes operator for automated pod stats collection, threshold-based alerting, and monitoring. Built with Go, Kubebuilder v4.6, and controller-runtime.

## Architecture

### System Overview

```mermaid
flowchart TB
    subgraph K8s["Kubernetes Cluster"]
        subgraph ONS["corium-operator-system"]
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
            GRAF[Grafana]
        end

        subgraph CRDs["Custom Resources"]
            CFG[JAXStatsConfig]
            COL[JAXStatsCollector]
            ALT[JAXStatsAlert]
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

### CRD Resource Hierarchy

```mermaid
graph LR
    CFG["JAXStatsConfig\n(global settings)"]
    COL["JAXStatsCollector\n(pod discovery)"]
    ALT["JAXStatsAlert\n(threshold rules)"]
    CM["ConfigMap\n(pod metrics)"]
    EV["K8s Events\n(alert notifications)"]

    CFG -->|referenced by| COL
    COL -->|referenced by| ALT
    COL -->|creates/owns| CM
    ALT -->|reads| CM
    ALT -->|emits| EV
```

### Collector Reconciliation Flow

```mermaid
sequenceDiagram
    participant CR as Collector CR Event
    participant R as Reconciler
    participant API as K8s API
    participant CM as ConfigMap

    CR->>R: Reconcile triggered
    R->>API: Get JAXStatsCollector
    R->>API: Get referenced JAXStatsConfig
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

```mermaid
flowchart TD
    START([Reconcile triggered]) --> FETCH[Fetch JAXStatsAlert]
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

### JAXStatsConfig

Global configuration for stats collection.

| Field | Type | Description |
|-------|------|-------------|
| `spec.enabled` | bool | Enable/disable collection |
| `spec.collectionInterval` | int32 | Collection interval in seconds (10-3600) |
| `spec.metrics` | []string | Metrics to collect |
| `spec.storageConfig.type` | enum | Storage backend: `prometheus`, `configmap`, `elasticsearch` |

### JAXStatsCollector

Discovers pods via label selectors and persists metrics to a ConfigMap.

| Field | Type | Description |
|-------|------|-------------|
| `spec.targetNamespace` | string | Namespace to discover pods in |
| `spec.selector` | LabelSelector | Pod label selector |
| `spec.configRef` | string | Name of referenced JAXStatsConfig |
| `status.discoveredPods` | int32 | Number of discovered pods |
| `status.metricsConfigMap` | string | Name of the managed ConfigMap |

### JAXStatsAlert

Evaluates threshold rules against collected metrics and emits K8s Events.

| Field | Type | Description |
|-------|------|-------------|
| `spec.enabled` | bool | Enable/disable alert evaluation |
| `spec.collectorRef` | string | Name of referenced JAXStatsCollector |
| `spec.cooldownPeriod` | string | Cooldown between alerts (e.g., "5m") |
| `spec.rules[].metric` | enum | `restart_count`, `not_ready_count`, `container_count`, `pod_count` |
| `spec.rules[].operator` | enum | `>`, `<`, `>=`, `<=`, `==` |
| `spec.rules[].threshold` | string | Numeric threshold value |
| `spec.rules[].severity` | enum | `critical`, `warning`, `info` |
| `status.firingAlertsCount` | int32 | Number of currently firing alerts |

## K8s Patterns Demonstrated

- **Finalizers** -- Config cleanup of dependent Collectors on deletion
- **Status Conditions** -- Standard `metav1.Condition` with `ObservedGeneration`
- **Owner References** -- Collector-owned ConfigMaps auto-deleted on CR removal
- **EventRecorder** -- AlertFiring/AlertResolved events visible via `kubectl get events`
- **Cross-resource reconciliation** -- Alert reads Collector's ConfigMap
- **Kubebuilder validation markers** -- Enum, MinLength, MinItems, Min/Max
- **Printer columns** -- `kubectl get` shows Enabled, Status, Pods, Firing counts
- **NetworkPolicies** -- Namespace isolation with explicit allow rules
- **Prometheus custom metrics** -- `jaxstats_discovered_pods`, `jaxstats_active_alerts`, `jaxstats_reconcile_errors_total`
- **Grafana dashboard** -- Pre-built JSON auto-provisioned via ConfigMap sidecar

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
kubectl get jsc,jscol,jsa
kubectl get configmap -l app.kubernetes.io/managed-by=jaxstats-operator
kubectl get events --field-selector reason=AlertFiring
```

## Full Demo

```bash
# One-command demo with KinD
cd .. && ./demo.sh
```

## License

MIT
