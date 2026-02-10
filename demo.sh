#!/usr/bin/env bash
set -euo pipefail

# Corium Operator Demo
# One-command demo: KinD cluster, operator, demo workloads, sample CRs

CLUSTER_NAME="corium-demo"
OPERATOR_IMG="corium-operator:demo"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=== Corium Operator Demo ==="
echo ""

# Check prerequisites
for cmd in docker kind kubectl; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "ERROR: $cmd is required but not installed."
    exit 1
  fi
done

# Clean up existing cluster if present
if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
  echo ">>> Deleting existing KinD cluster..."
  kind delete cluster --name "$CLUSTER_NAME"
fi

# Create KinD cluster
echo ">>> Creating KinD cluster..."
kind create cluster --name "$CLUSTER_NAME" --config "$SCRIPT_DIR/kind-config.yaml" --wait 60s

echo ">>> Cluster ready."
kubectl cluster-info --context "kind-${CLUSTER_NAME}"
echo ""

# Build operator image
echo ">>> Building operator image..."
cd "$SCRIPT_DIR/operator"
docker build -t "$OPERATOR_IMG" -f Dockerfile .
kind load docker-image "$OPERATOR_IMG" --name "$CLUSTER_NAME"

# Install CRDs
echo ">>> Installing CRDs..."
make install

# Deploy namespaces and network policies
echo ">>> Creating namespaces and network policies..."
kubectl apply -f "$SCRIPT_DIR/k8s/namespaces.yaml"
kubectl apply -f "$SCRIPT_DIR/k8s/network-policies.yaml"

# Run operator locally in background (for demo simplicity)
echo ">>> Starting operator (background)..."
make run &
OPERATOR_PID=$!

# Give operator time to start
sleep 5

# Deploy demo workload
echo ">>> Deploying demo workload (3x nginx pods)..."
kubectl apply -f config/samples/demo-workload.yaml

# Wait for pods to be ready
echo ">>> Waiting for demo pods..."
kubectl wait --for=condition=ready pod -l app=demo-workload --timeout=60s || true

# Apply sample CRs
echo ">>> Applying sample JAXStats resources..."
kubectl apply -f config/samples/stats_v1alpha1_jaxstatsconfig.yaml
kubectl apply -f config/samples/stats_v1alpha1_jaxstatscollector.yaml
kubectl apply -f config/samples/stats_v1alpha1_jaxstatsalert.yaml

# Wait for reconciliation
echo ">>> Waiting for reconciliation (10s)..."
sleep 10

# Show results
echo ""
echo "=== Results ==="
echo ""

echo "--- JAXStatsConfig ---"
kubectl get jaxstatsconfigs -o wide 2>/dev/null || kubectl get jaxstatsconfigs
echo ""

echo "--- JAXStatsCollector ---"
kubectl get jaxstatscollectors -o wide 2>/dev/null || kubectl get jaxstatscollectors
echo ""

echo "--- JAXStatsAlert ---"
kubectl get jaxstatsalerts -o wide 2>/dev/null || kubectl get jaxstatsalerts
echo ""

echo "--- Metrics ConfigMap ---"
kubectl get configmap -l app.kubernetes.io/managed-by=jaxstats-operator
echo ""

CMNAME=$(kubectl get configmap -l app.kubernetes.io/managed-by=jaxstats-operator -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
if [ -n "$CMNAME" ]; then
  echo "--- ConfigMap Contents (metrics.json) ---"
  kubectl get configmap "$CMNAME" -o jsonpath='{.data.metrics\.json}' | python3 -m json.tool 2>/dev/null || \
    kubectl get configmap "$CMNAME" -o jsonpath='{.data.metrics\.json}'
  echo ""
fi

echo "--- K8s Events (alerts) ---"
kubectl get events --field-selector reason=AlertFiring 2>/dev/null || echo "(no alerts firing)"
kubectl get events --field-selector reason=AlertResolved 2>/dev/null || echo "(no alerts resolved)"
echo ""

echo "--- Network Policies ---"
kubectl get networkpolicies -A
echo ""

echo "=== Demo Complete ==="
echo ""
echo "Operator is running in background (PID: $OPERATOR_PID)"
echo "To stop: kill $OPERATOR_PID"
echo "To clean up: kind delete cluster --name $CLUSTER_NAME"
echo ""
echo "Try:"
echo "  kubectl get jsc,jscol,jsa"
echo "  kubectl describe jaxstatsalert demo-alert"
echo "  kubectl get configmap demo-collector-metrics -o yaml"
echo "  kubectl get events --field-selector reason=AlertFiring"
