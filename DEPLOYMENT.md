# Corium Deployment Guide

This guide outlines the steps to deploy Corium to a Kubernetes cluster.

## Prerequisites

- Kubernetes cluster (v1.19 or later)
- kubectl configured to communicate with your cluster
- Docker installed and configured
- Access to a container registry (optional, but recommended)

## Building the Docker Image

1. Build the Docker image:
   ```bash
   docker build -t corium:latest .
   ```

2. (Optional) Push to a container registry:
   ```bash
   docker tag corium:latest your-registry/corium:latest
   docker push your-registry/corium:latest
   ```

## Deploying to Kubernetes

1. Deploy the application:
   ```bash
   kubectl apply -f k8s/deployment.yaml
   kubectl apply -f k8s/service.yaml
   kubectl apply -f k8s/ingress.yaml
   ```

2. Verify the deployment:
   ```bash
   kubectl get deployments
   kubectl get pods
   kubectl get services
   kubectl get ingress
   ```

## Configuration

### Ingress Configuration

Before applying the ingress configuration:

1. Update `k8s/ingress.yaml`:
   - Replace `corium.example.com` with your actual domain
   - Update the TLS secret name if needed

2. Create TLS secret (if using HTTPS):
   ```bash
   kubectl create secret tls corium-tls --cert=path/to/cert.pem --key=path/to/key.pem
   ```

### Resource Configuration

You can adjust resource limits and requests in `k8s/deployment.yaml`:
- CPU and memory limits
- Number of replicas
- Health check parameters

## Troubleshooting

1. Check pod logs:
   ```bash
   kubectl logs -f $(kubectl get pods -l app=corium -o jsonpath="{.items[0].metadata.name}")
   ```

2. Check pod status:
   ```bash
   kubectl describe pod $(kubectl get pods -l app=corium -o jsonpath="{.items[0].metadata.name}")
   ```

3. Check ingress status:
   ```bash
   kubectl describe ingress corium
   ```

## Cleanup

To remove the deployment:
```bash
kubectl delete -f k8s/
```

## Additional Notes

- Make sure your cluster has an ingress controller installed
- Configure DNS to point to your cluster's ingress IP
- Consider setting up monitoring and logging solutions
- Review security best practices for production deployments

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
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

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

## Using Local Images with KinD (Development)

By default, the deployment uses a remote image from your registry.
For local development with KinD, you can patch the deployment to use your local image and prevent Kubernetes from trying to pull it from a registry:

```bash
# Build your image locally
docker build -t corium:latest .

# Load it into your KinD cluster
kind load docker-image corium:latest --name corium-cluster

# Patch the deployment to use the local image and never pull from a registry
kubectl set image deployment/corium corium=corium:latest
kubectl patch deployment corium -p '{"spec":{"template":{"spec":{"containers":[{"name":"corium","imagePullPolicy":"Never"}]}}}}'
```

To revert to using a remote image (e.g., for production), update your deployment manifest to use your registry and the default `IfNotPresent` policy, then re-apply:

```yaml
image: your-registry/corium:latest
imagePullPolicy: IfNotPresent
```

## Production Deployment

### 1. Prepare Environment

1. Create a `.env` file with your configuration:
   ```env
   NODE_ENV=production
   ```

2. Build the production image:
   ```bash
   docker build -t your-registry/corium:latest .
   docker tag corium:latest your-registry/corium:latest
   docker push your-registry/corium:latest
   ```

### 2. Deploy to Kubernetes

```bash
kubectl apply -f k8s/
```

### 3. Verify deployment:
```bash
kubectl get pods
kubectl get services
```

## Monitoring and Maintenance

### Health Checks

```bash
kubectl get pods
kubectl logs -f deployment/corium
```

### Scaling

```bash
kubectl scale deployment corium --replicas=3
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
kubectl get all -l app=corium -o yaml > corium-backup.yaml
```

### Recovery

```bash
kubectl apply -f corium-backup.yaml
```

## Additional Resources

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [KinD Documentation](https://kind.sigs.k8s.io/)
- [Docker Documentation](https://docs.docker.com/) 