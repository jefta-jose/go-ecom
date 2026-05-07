# postgres-db Helm Chart

This chart deploys a minimal PostgreSQL instance in the `ecom-go-db` namespace, with access restricted to the API service in the `ecom-go-api` namespace.

## Usage

1. **Create the namespace:**
   ```sh
   kubectl create namespace ecom-go-db
   ```
2. **Install the chart:**
   ```sh
   helm install postgres-db ./helm/postgres-db -n ecom-go-db
   ```

## Access from Another Namespace

- The PostgreSQL service will be available at:
  ```
  postgres-db.ecom-go-db.svc.cluster.local:5432
  ```
- Only pods with label `app: ecom-go-api` in the `ecom-go-api` namespace can connect (enforced by NetworkPolicy).

## NetworkPolicy
- See `templates/networkpolicy.yaml` for details.

## Environment Variables
- Set in `values.yaml`:
  - `POSTGRES_USER`
  - `POSTGRES_PASSWORD`
  - `POSTGRES_DB`

## Persistent Storage
- PVC is created for `/var/lib/postgresql/data`.

# If you want to run tests against this db locally you will have to port-forward the connection
run 
```
kubectl port-forward svc/postgres-db -n ecom-go-db 5432:5432
```

this way you test the db directly without depending on the deployed api in helm