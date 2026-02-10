# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Corium is a Kubernetes-native platform with two main components:
1. **Main App** — A Next.js 14 portfolio/landing page (TypeScript, Tailwind CSS, Framer Motion)
2. **Operator** — A Kubernetes operator for managing JaxStats resources (Go, Kubebuilder v4.6.0)

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
```

### Kubernetes Operator (operator/)
```bash
make manifests           # Generate CRDs
make generate            # Generate DeepCopy methods
make build               # Build manager binary
make run                 # Run controller locally
make test                # Run unit tests (uses envtest)
make test-e2e            # Run e2e tests (requires KinD)
make lint                # golangci-lint
make lint-fix            # golangci-lint with auto-fix
make install             # Install CRDs into cluster
make deploy IMG=<image>  # Deploy controller to cluster
```

### Docker & KinD
```bash
docker build -t corium:latest .                          # Build main app image
kind create cluster --name corium-cluster --config kind-config.yaml  # Local cluster
kind load docker-image corium:latest --name corium-cluster
kubectl apply -f k8s/                                    # Deploy all manifests
```

## Architecture

### Main App (src/)
- **Next.js 14 App Router** with `src/app/` directory structure
- Client components use `"use client"` directive for Framer Motion animations
- Path alias: `@/*` maps to `./src/*` (configured in tsconfig.json)
- Components: `src/components/` (Hero, Experience, Education, Skills, Projects)
- Standalone output mode; TypeScript build errors ignored in next.config.js

### Kubernetes Operator (operator/)
- Kubebuilder v4.6.0, API group `stats.corium.io`, version `v1alpha1`
- Three CRDs: `JAXStatsConfig`, `JAXStatsCollector`, `JAXStatsAlert`
- Types: `operator/api/v1alpha1/`
- Controllers: `operator/internal/controller/`
- Kustomize manifests: `operator/config/`

### Kubernetes Manifests (k8s/)
- Corium deployment (3 replicas)
- Services, Ingress, Namespace

## Git Conventions

- Do NOT add a `Co-Authored-By` line to commit messages.

## CI/CD

GitHub Actions pipeline (`.github/workflows/main.yml`): test → build → deploy (Docker Hub push on main branch only). Requires `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` secrets.
