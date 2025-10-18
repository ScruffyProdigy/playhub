# PlayHub Environment Configuration

This document explains how environment configuration works in PlayHub, particularly for the frontend application across different deployment environments.

## Overview

PlayHub uses a Docker-based approach to inject environment variables into the frontend application at runtime. This allows the same Docker image to be deployed to different environments (local, staging, production) with different configurations.

## How It Works

### 1. Docker Entrypoint Script

The frontend Dockerfile includes a script that runs when the container starts:

```dockerfile
# Create a script to inject environment variables
RUN echo '#!/bin/sh' > /docker-entrypoint.d/30-env-injection.sh && \
    echo 'echo "window.env = {" > /usr/share/nginx/html/env.js' >> /docker-entrypoint.d/30-env-injection.sh && \
    echo 'echo "  REACT_APP_ENV: \"$REACT_APP_ENV\"," >> /usr/share/nginx/html/env.js' >> /docker-entrypoint.d/30-env-injection.sh && \
    echo 'echo "  REACT_APP_API_BASE_URL: \"$REACT_APP_API_BASE_URL\"" >> /usr/share/nginx/html/env.js' >> /docker-entrypoint.d/30-env-injection.sh && \
    echo 'echo "};" >> /usr/share/nginx/html/env.js' >> /docker-entrypoint.d/30-env-injection.sh && \
    chmod +x /docker-entrypoint.d/30-env-injection.sh
```

This script creates an `env.js` file in the nginx document root with the current environment variables.

### 2. Frontend Loading

The frontend loads this configuration in `index.html`:

```html
<script src="/env.js" type="module"></script>
```

The frontend can then access these variables via `window.env`:

```javascript
const apiBaseUrl = window.env.REACT_APP_API_BASE_URL;
const environment = window.env.REACT_APP_ENV;
```

## Environment Configurations

### Local Development

**File**: `k8s/env/local.yaml`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: lobby-frontend-config
  namespace: playhub
  labels: { env: local }
data:
  REACT_APP_ENV: local
  REACT_APP_API_BASE_URL: "http://localhost:8081"
```

**Deployment**: `./deploy-local.sh`

- Uses `minikube` context
- Backend accessible via port-forward on `localhost:8081`
- Frontend accessible via port-forward on `localhost:8080`

### Staging

**File**: `k8s/env/staging.yaml`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: lobby-frontend-config
  namespace: playhub-staging
  labels: { env: staging }
data:
  REACT_APP_ENV: staging
  REACT_APP_API_BASE_URL: "https://api-staging.playhub.com"
```

**Deployment**: `./deploy-staging.sh`

- Uses `staging-cluster` context
- 2 replicas for high availability
- Uses `staging` tagged Docker images
- Higher resource limits

### Production

**File**: `k8s/env/production.yaml`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: lobby-frontend-config
  namespace: playhub-production
  labels: { env: production }
data:
  REACT_APP_ENV: production
  REACT_APP_API_BASE_URL: "https://api.playhub.com"
```

**Deployment**: `./deploy-production.sh`

- Uses `production-cluster` context
- 3 replicas for high availability
- Uses `production` tagged Docker images
- Highest resource limits
- Longer health check delays

## Deployment Scripts

### Local Deployment

```bash
./deploy-local.sh
```

This script:
1. Sets kubectl context to `minikube`
2. Creates the `playhub` namespace
3. Applies base configurations (backend, ingress)
4. Applies local environment configuration
5. Waits for deployments to be ready
6. Shows deployment status and access instructions

### Staging Deployment

```bash
./deploy-staging.sh
```

This script:
1. Sets kubectl context to `staging-cluster`
2. Creates the `playhub-staging` namespace
3. Applies base configurations
4. Applies staging environment configuration
5. Waits for deployments to be ready
6. Shows deployment status

### Production Deployment

```bash
./deploy-production.sh
```

This script:
1. Sets kubectl context to `production-cluster`
2. Creates the `playhub-production` namespace
3. Applies base configurations
4. Applies production environment configuration
5. Waits for deployments to be ready
6. Shows deployment status

## Environment Variables

### Available Variables

- `REACT_APP_ENV`: Environment identifier (`local`, `staging`, `production`)
- `REACT_APP_API_BASE_URL`: Base URL for the GraphQL API

### Adding New Variables

To add new environment variables:

1. **Update the Dockerfile script** to include the new variable:
   ```dockerfile
   echo 'echo "  REACT_APP_NEW_VAR: \"$REACT_APP_NEW_VAR\"," >> /usr/share/nginx/html/env.js' >> /docker-entrypoint.d/30-env-injection.sh
   ```

2. **Add the variable to all environment ConfigMaps**:
   ```yaml
   data:
     REACT_APP_ENV: local
     REACT_APP_API_BASE_URL: "http://localhost:8081"
     REACT_APP_NEW_VAR: "local-value"
   ```

3. **Update the frontend code** to use the new variable:
   ```javascript
   const newVar = window.env.REACT_APP_NEW_VAR;
   ```

## Benefits

1. **Single Docker Image**: The same image works in all environments
2. **Runtime Configuration**: No need to rebuild for different environments
3. **Kubernetes Native**: Uses ConfigMaps for environment-specific values
4. **Secure**: Sensitive values can be stored in Secrets and referenced in ConfigMaps
5. **Easy Deployment**: Simple scripts for each environment

## Troubleshooting

### Check Environment Configuration

```bash
# View the current ConfigMap
kubectl get configmap lobby-frontend-config -n playhub -o yaml

# Check if env.js is being generated correctly
kubectl exec -n playhub deployment/lobby-frontend -- cat /usr/share/nginx/html/env.js
```

### Verify Frontend Access

```bash
# Check if the frontend is serving env.js
curl http://localhost:8080/env.js

# Check the generated content
kubectl port-forward -n playhub svc/lobby-frontend 8080:80 &
curl http://localhost:8080/env.js
```

### Common Issues

1. **env.js not found**: Check if the Docker entrypoint script is running
2. **Wrong API URL**: Verify the ConfigMap values match your environment
3. **CORS errors**: Ensure the API URL is accessible from the frontend domain
