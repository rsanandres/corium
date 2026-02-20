# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Corium is a Kubernetes-native platform with two main components:
1. **Main App** — A Next.js 14 portfolio/landing page (TypeScript, Tailwind CSS, Framer Motion)
2. **Operator** — A Kubernetes operator for monitoring workloads (Go 1.24, Kubebuilder v4.6.0)

## Build & Development Commands

### Main App (root directory)
```bash
npm install              # Install dependencies
npm run dev              # Dev server on :3000
npm run build            # Production build
npm run lint             # ESLint
npm run type-check       # TypeScript checking (tsc --noEmit)
npm test                 # Jest tests
npm run test:watch       # Jest watch mode
npm run test:coverage    # Jest with coverage
npm test -- Foo.test.tsx                  # Run a single test file
npm test -- --testNamePattern="some test" # Run tests matching a name
```

### Kubernetes Operator (operator/)
```bash
make manifests           # Generate CRDs
make generate            # Generate DeepCopy methods
make build               # Build manager binary
make fmt                 # Format Go code
make run                 # Run controller locally
make test                # Run unit tests (uses envtest)
make test-e2e            # Run e2e tests (requires KinD)
make lint                # golangci-lint
make lint-fix            # golangci-lint with auto-fix
make install             # Install CRDs into cluster
make deploy IMG=<image>  # Deploy controller to cluster
make build-installer     # Generate dist/install.yaml
```

Running a single operator test (Ginkgo pattern match):
```bash
cd operator && go test ./internal/controller -run TestControllers -ginkgo.focus "pattern"
```

### Docker & KinD
```bash
docker build -t corium:latest .
kind create cluster --name corium-cluster --config kind-config.yaml
kind load docker-image corium:latest --name corium-cluster
kubectl apply -f k8s/
```

## Architecture

### Main App (src/)
- **Next.js 14 App Router** with `src/app/` directory structure
- Client components use `"use client"` directive for Framer Motion animations
- Path alias: `@/*` maps to `./src/*` (configured in tsconfig.json)
- Components: `src/components/` (Hero, Experience, Education, Skills, Projects)
- Standalone output mode; TypeScript build errors ignored in next.config.js
- Testing: Jest 29 + @testing-library/react (jsdom environment). No test files exist yet.

### Kubernetes Operator (operator/)
- **API group** `monitor.corium.io`, version `v1alpha1`
- **Three CRDs with a dependency chain:**
  - `CoriumMonitorConfig` — Base configuration (namespace targeting, collection params)
  - `CoriumMonitorCollector` — References a Config via `spec.configRef`; collects pod metrics and writes them to an owned ConfigMap (`{name}-metrics`)
  - `CoriumMonitorAlert` — References a Collector via `spec.collectorRef`; evaluates alert rules against collected metrics
- Types: `operator/api/v1alpha1/`
- Controllers: `operator/internal/controller/`
- Pure logic extracted into `metrics_collector.go` and `alert_evaluator.go`
- Custom Prometheus metrics registered in `metrics.go` (`corium_discovered_pods`, `corium_active_alerts`, `corium_reconcile_errors_total`)
- Config controller uses a finalizer (`monitor.corium.io/config-cleanup`) for cascading cleanup of dependent Collectors
- Kustomize manifests: `operator/config/`

### Operator Testing
- **Framework:** Ginkgo v2 + Gomega (BDD-style), NOT standard Go `testing`
- `suite_test.go` bootstraps envtest (in-memory K8s API server) shared across all controller tests
- Tests create CRs in `BeforeEach`, reconcile directly via `reconciler.Reconcile()`, and clean up in `AfterEach`
- Uses `record.FakeRecorder` for asserting K8s events
- envtest binaries auto-discovered from `operator/bin/k8s/`

### Kubernetes Manifests (k8s/)
- **Three namespaces:** `operator-system` (controllers/CRDs), `corium-monitoring` (Prometheus + Grafana), `corium-workloads` (demo pods)
- Network policies: default-deny per namespace with explicit allow rules
- Monitoring stack: `k8s/monitoring/` contains Prometheus Helm values, Grafana dashboard JSON, and ServiceMonitor config
- Main app deployment: 3 replicas, `imagePullPolicy: Never` for local KinD dev

## Git Conventions

- Do NOT add a `Co-Authored-By` line to commit messages.

## CI/CD

GitHub Actions pipeline (`.github/workflows/main.yml`): test → operator-test → build → deploy (Docker Hub push on main branch only). Uses Node 20 and Go 1.24. Requires `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` secrets.

