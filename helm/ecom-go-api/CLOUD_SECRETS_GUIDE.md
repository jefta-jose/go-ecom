# Cloud Secret Manager Integration Guide

This Helm chart supports fetching secrets from cloud secret managers using the External Secrets Operator.

## Overview

Instead of hardcoding sensitive data (like database connection strings) in your Helm chart, you can store them securely in:
- **AWS Secrets Manager**
- **Azure Key Vault**
- **GCP Secret Manager**

The External Secrets Operator syncs these cloud secrets into Kubernetes secrets automatically.

---

## Prerequisites

### 1. Install External Secrets Operator

```bash
# Add the External Secrets Operator Helm repo
helm repo add external-secrets https://charts.external-secrets.io
helm repo update

# Install the operator in its own namespace
helm install external-secrets external-secrets/external-secrets \
  -n external-secrets-system \
  --create-namespace
```

### 2. Configure Cloud Provider Access

Choose your cloud provider and follow the setup instructions below.

---

## AWS Secrets Manager Setup

### Option A: Using IRSA (Recommended for EKS)

**1. Create IAM Policy:**
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue",
        "secretsmanager:DescribeSecret"
      ],
      "Resource": "arn:aws:secretsmanager:REGION:ACCOUNT_ID:secret:ecom-go-api/*"
    }
  ]
}
```

**2. Create IAM Role with Trust Policy for EKS:**
```bash
eksctl create iamserviceaccount \
  --name ecom-go-api \
  --namespace default \
  --cluster YOUR_CLUSTER_NAME \
  --attach-policy-arn arn:aws:iam::ACCOUNT_ID:policy/EcomGoApiSecretsPolicy \
  --approve
```

**3. Update values.yaml:**
```yaml
externalSecrets:
  enabled: true
  provider: aws
  aws:
    region: us-east-1
    secretName: ecom-go-api/db-string
    secretKey: DBSTRING
    roleArn: arn:aws:iam::ACCOUNT_ID:role/EcomGoApiSecretsRole

serviceAccount:
  create: true
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::ACCOUNT_ID:role/EcomGoApiSecretsRole
```

---

## Azure Key Vault Setup

### Option A: Using Managed Identity (Recommended for AKS)

**1. Enable Workload Identity on AKS:**
```bash
az aks update \
  --resource-group YOUR_RESOURCE_GROUP \
  --name YOUR_CLUSTER_NAME \
  --enable-workload-identity \
  --enable-oidc-issuer
```

**2. Create Managed Identity and Grant Access:**
```bash
# Create user-assigned managed identity
az identity create \
  --name ecom-go-api-identity \
  --resource-group YOUR_RESOURCE_GROUP

# Get the identity client ID
CLIENT_ID=$(az identity show --name ecom-go-api-identity --resource-group YOUR_RESOURCE_GROUP --query clientId -o tsv)

# Grant Key Vault access
az keyvault set-policy \
  --name YOUR_KEY_VAULT_NAME \
  --spn $CLIENT_ID \
  --secret-permissions get list
```

**3. Update values.yaml:**
```yaml
externalSecrets:
  enabled: true
  provider: azure
  azure:
    vaultUrl: https://YOUR_KEY_VAULT_NAME.vault.azure.net/
    secretName: db-connection-string
    useManagedIdentity: true

serviceAccount:
  create: true
  annotations:
    azure.workload.identity/client-id: YOUR_CLIENT_ID
```

### Option B: Using Service Principal

**1. Create Service Principal and Grant Access:**
```bash
# Create service principal
az ad sp create-for-rbac --name ecom-go-api-sp

# Grant Key Vault access
az keyvault set-policy \
  --name YOUR_KEY_VAULT_NAME \
  --spn YOUR_CLIENT_ID \
  --secret-permissions get list
```

**2. Create Kubernetes Secret:**
```bash
kubectl create secret generic azure-credentials \
  --from-literal=client-id=YOUR_CLIENT_ID \
  --from-literal=client-secret=YOUR_CLIENT_SECRET
```

**3. Update values.yaml:**
```yaml
externalSecrets:
  enabled: true
  provider: azure
  azure:
    vaultUrl: https://YOUR_KEY_VAULT_NAME.vault.azure.net/
    secretName: db-connection-string
    useManagedIdentity: false
    tenantId: YOUR_TENANT_ID
    credentialsSecretName: azure-credentials
```

**4. Create Secret in Key Vault:**
```bash
az keyvault secret set \
  --vault-name YOUR_KEY_VAULT_NAME \
  --name db-connection-string \
  --value "host=postgres user=postgres password=xxx dbname=go_ecom sslmode=require"
```

---

## GCP Secret Manager Setup

### Option A: Using Workload Identity (Recommended for GKE)

**1. Enable Workload Identity on GKE:**
```bash
gcloud container clusters update YOUR_CLUSTER_NAME \
  --workload-pool=YOUR_PROJECT_ID.svc.id.goog
```

**2. Create Service Account and Grant Permissions:**
```bash
# Create GCP service account
gcloud iam service-accounts create ecom-go-api-sa \
  --display-name="Ecom Go API Service Account"

# Grant Secret Manager access
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:ecom-go-api-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"

# Bind Kubernetes service account to GCP service account
gcloud iam service-accounts add-iam-policy-binding \
  ecom-go-api-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com \
  --role roles/iam.workloadIdentityUser \
  --member "serviceAccount:YOUR_PROJECT_ID.svc.id.goog[default/ecom-go-api]"
```

**3. Update values.yaml:**
```yaml
externalSecrets:
  enabled: true
  provider: gcp
  gcp:
    projectId: YOUR_PROJECT_ID
    secretName: db-connection-string
    secretVersion: latest
    useWorkloadIdentity: true
    clusterLocation: us-central1
    clusterName: YOUR_CLUSTER_NAME

serviceAccount:
  create: true
  annotations:
    iam.gke.io/gcp-service-account: ecom-go-api-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

---

## Deployment

After configuring your chosen cloud provider:

```bash
# Deploy the Helm chart
helm install ecom-go-api ./helm/ecom-go-api

# Verify the ExternalSecret is syncing
kubectl get externalsecret
kubectl describe externalsecret ecom-go-api-external-secret

# Verify the Kubernetes secret was created
kubectl get secret ecom-go-api-db-secret
```

---

## Local Development

For local development, keep `externalSecrets.enabled: false` and use the local secret:

```yaml
db:
  dsn: "host=localhost user=postgres password=postgres123 dbname=go_ecom sslmode=disable"

externalSecrets:
  enabled: false  # Uses local db.dsn value instead
```

---