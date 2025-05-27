# Corium Deployment Guide

This guide provides detailed instructions for deploying Corium to various Kubernetes environments.

## Prerequisites

- Kubernetes cluster (local or remote)
- kubectl configured and pointing to your cluster
- Docker installed locally
- Node.js and npm installed

## Recommended Tools

### Kubernetes Management Tools

For better Kubernetes cluster management and monitoring, we recommend installing these tools:

1. **k9s** - Terminal-based UI for Kubernetes
   ```bash
   # macOS
   brew install k9s
   
   # Linux
   curl -sS https://webinstall.dev/k9s | bash
   ```

2. **Lens** - Kubernetes IDE
   - Download from [Lens Official Website](https://k8slens.dev/)
   - Available for macOS, Windows, and Linux
   - Provides a powerful GUI for cluster management

These tools will help you:
- Monitor cluster resources in real-time
- Debug deployments and pods
- Manage configurations
- View logs and events
- Scale resources efficiently

## Local Development with KinD

### 1. Set up KinD Cluster

```bash
# Install KinD if not already installed
# For macOS:
brew install kind

# For Linux:
# curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
# chmod +x ./kind
# sudo mv ./kind /usr/local/bin/kind

# Create a new cluster
kind create cluster --name corium-cluster
```

### 2. Build and Deploy

```bash
# Build the Docker image
docker build -t corium:latest .

# Load the image into KinD
kind load docker-image corium:latest --name corium-cluster

# Apply Kubernetes manifests
kubectl apply -f k8s/
```

## Production Deployment

### 1. Prepare Environment

1. Create a `.env` file with your configuration:
   ```env
   NODE_ENV=production
   KUBERNETES_NAMESPACE=corium
   ```

2. Build the production image:
   ```bash
   # Build the Docker image
   docker build -t your-registry/corium:latest .
   
   # If using a private registry, tag and push the image
   docker tag corium:latest your-registry/corium:latest
   docker push your-registry/corium:latest
   ```

   The Dockerfile uses a multi-stage build process:
   - First stage builds the TypeScript application and handles submodules
   - Second stage creates a minimal production image
   - Includes proper handling of the JaxStats submodule

### 2. Deploy to Kubernetes

1. Create namespace:
   ```bash
   kubectl create namespace corium
   ```

2. Apply Kubernetes manifests:
   ```bash
   kubectl apply -f k8s/
   ```

3. Verify deployment:
   ```bash
   kubectl get pods -n corium
   kubectl get services -n corium
   ```

## Kubernetes Manifests

The `k8s/` directory contains the following manifests:

- `deployment.yaml` - Main application deployment
- `service.yaml` - Service definitions
- `ingress.yaml` - Ingress configuration
- `configmap.yaml` - Configuration management
- `secrets.yaml` - Sensitive data management

## Monitoring and Maintenance

### Health Checks

```bash
# Check pod status
kubectl get pods -n corium

# View logs
kubectl logs -f deployment/corium -n corium
```

### Scaling

```bash
# Scale deployment
kubectl scale deployment corium --replicas=3 -n corium
```

## Troubleshooting

### Common Issues

1. **Image Pull Errors**
   - Verify image exists in registry
   - Check image pull secrets

2. **Pod CrashLoopBackOff**
   - Check pod logs
   - Verify environment variables
   - Check resource limits

3. **Service Unavailable**
   - Verify service endpoints
   - Check network policies
   - Verify ingress configuration

## Security Considerations

- Use Kubernetes secrets for sensitive data
- Implement network policies
- Enable RBAC
- Regular security updates
- Monitor for vulnerabilities

## Backup and Recovery

### Backup

```bash
# Backup Kubernetes resources
kubectl get all -n corium -o yaml > corium-backup.yaml
```

### Recovery

```bash
# Restore from backup
kubectl apply -f corium-backup.yaml
```

## Additional Resources

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [KinD Documentation](https://kind.sigs.k8s.io/)
- [Docker Documentation](https://docs.docker.com/) 