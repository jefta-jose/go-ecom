# build the api image
docker build -t ecom-go-api-project:latest .

# ecom-go-api Helm Chart

This Helm chart deploys the ecom-go-api application on Kubernetes. It is designed for production and local development, supporting best practices for configuration, scaling, and cloud secret management.

## Features
- Deploys a Go API app with configurable replicas, resources, and probes
- Supports Ingress (with NGINX or other controllers)
- Cloud-agnostic secret management via External Secrets Operator (AWS, Azure, GCP)
- Horizontal Pod Autoscaler (HPA) support
- ServiceAccount for cloud IAM integration
- Well-commented values.yaml for easy customization

## Quick Start

### 1. Prerequisites
- Kubernetes cluster (local or cloud)
- Helm 3.x installed
- (Optional) External Secrets Operator for cloud secret integration

### 2. Install the Chart
```sh
# Create a namespace (recommended)
kubectl create namespace ecom-go-api

# Install the chart
helm install ecom-go-api ./helm/ecom-go-api --namespace ecom-go-api
```

### 3. Upgrade the Chart
```sh
# Change values in values.yaml or use --set
helm upgrade ecom-go-api ./helm/ecom-go-api --namespace ecom-go-api
```

### 4. Check Status
```sh
helm status ecom-go-api --namespace ecom-go-api
kubectl get all -n ecom-go-api
```

## Configuration
Edit `values.yaml` to customize:
- `replicaCount`: Number of pods
- `image`: Docker image settings
- `service`: Service type and ports
- `ingress`: Enable and configure Ingress
- `resources`: CPU/memory requests and limits
- `db.dsn`: Database connection string (for local dev)
- `externalSecrets`: Enable and configure cloud secret manager integration
- `autoscaling`: Enable and tune HPA

See comments in `values.yaml` for detailed guidance.

## Cloud Secret Management
- See `CLOUD_SECRETS_GUIDE.md` for step-by-step instructions on using AWS Secrets Manager, Azure Key Vault, or GCP Secret Manager with this chart.

## Debugging & Troubleshooting
- For common Kubernetes issues and conversational debugging, see `.claude/skills/k8-debugger/SKILL.md`.
- Use `kubectl describe` and `kubectl logs` to investigate pod and service issues.

## Uninstall
```sh
helm uninstall ecom-go-api --namespace ecom-go-api
```

## Resources
- [Helm documentation](https://helm.sh/docs/)
- [External Secrets Operator](https://external-secrets.io)
- [Kubernetes documentation](https://kubernetes.io/docs/)

# confirm image has laoded
# understand that nodes inside a kind cluster are just docker containers running kindest-node
docker exec -it kind-control-plane
crictl images

# 3️⃣ Update database DSN
# Edit values.yaml and set `db.dsn` to your PostgreSQL connection string

# 4️⃣ Install Helm chart
helm install ecom-go-api ./helm/ecom-go-api -n <namespace-of-your-choice>

# 5️⃣ Access the API via port-forward OR ingress

# Option A: Direct port-forward to service (simplest)
kubectl port-forward svc/ecom-go-api-go-api-chart 8080:80
curl http://localhost:8080/health

# Option B: Via ingress (requires nginx ingress controller)
# First, map the hostname (127.0.0.1 = localhost)
echo "127.0.0.1 go-api.local" | sudo tee -a /etc/hosts

# Then access via hostname
curl http://go-api.local/health
```

### Verify Deployment

```bash
# List all resources for this release
kubectl get all -l app.kubernetes.io/instance=ecom-api

# Check pod status
kubectl get pods

# Inspect pod logs (tail last 50 lines)
kubectl logs -l app.kubernetes.io/instance=ecom-api --tail=50

# Describe a pod if debugging is needed
kubectl describe pod <pod-name>
```

---