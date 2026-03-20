# Health API Documentation

## Overview

The Health API provides two endpoints for monitoring the Ontology Bot backend service:

1. **Liveness Probe** (`/ontology_bot/v1/health`) - Indicates if the service is running
2. **Readiness Probe** (`/ontology_bot/v1/health/ready`) - Indicates if the service is ready to handle requests

## Endpoints

### 1. Health Check (Liveness)

**Endpoint**: `GET /ontology_bot/v1/health`

**Purpose**: Returns the current health status of the service. This endpoint is typically used by load balancers to determine if the service is alive.

**Response**:
```json
{
  "status": "healthy",
  "timestamp": "2026-03-19T10:30:00Z",
  "version": "1.0.0"
}
```

**Fields**:
- `status`: Service status (always "healthy" if endpoint responds)
- `timestamp`: Server timestamp when check was performed
- `version`: Application version

**Example**:
```bash
curl http://localhost:4441/ontology_bot/v1/health
```

### 2. Readiness Check

**Endpoint**: `GET /ontology_bot/v1/health/ready`

**Purpose**: Indicates if the service is ready to accept requests. This checks dependencies like database connectivity and external services.

**Response**:
```json
{
  "ready": true,
  "timestamp": "2026-03-19T10:30:00Z"
}
```

**Fields**:
- `ready`: Boolean indicating if service can handle requests
- `timestamp`: Server timestamp when check was performed

**Example**:
```bash
curl http://localhost:4441/ontology_bot/v1/health/ready
```

## Integration Examples

### Kubernetes Health Probes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ontology-bot
spec:
  template:
    spec:
      containers:
      - name: ontology-bot
        livenessProbe:
          httpGet:
            path: /ontology_bot/v1/health
            port: 4441
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ontology_bot/v1/health/ready
            port: 4441
          initialDelaySeconds: 5
          periodSeconds: 5
```

### Docker Health Check

```dockerfile
FROM golang:1.24-alpine AS builder
# ... build steps ...

FROM alpine:latest
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:4441/ontology_bot/v1/health || exit 1
CMD ["./server"]
```

### Monitoring Scripts

```bash
#!/bin/bash
# Simple health check script

HEALTH_URL="http://localhost:4441/ontology_bot/v1/health"
READY_URL="http://localhost:4441/ontology_bot/v1/health/ready"

check_health() {
  response=$(curl -s -o /dev/null -w "%{http_code}" "$HEALTH_URL")
  if [ "$response" = "200" ]; then
    echo "✓ Service is healthy"
  else
    echo "✗ Service health check failed (HTTP $response)"
    exit 1
  fi
}

check_ready() {
  response=$(curl -s "$READY_URL")
  ready=$(echo "$response" | grep -o '"ready":[^,}]*' | grep -o '[^:]*$')
  if [ "$ready" = "true" ]; then
    echo "✓ Service is ready"
  else
    echo "✗ Service is not ready"
    exit 1
  fi
}

check_health
check_ready
```

## Testing the Health API

### Using curl

```bash
# Test liveness
curl -v http://localhost:4441/ontology_bot/v1/health

# Test readiness
curl -v http://localhost:4441/ontology_bot/v1/health/ready
```

### Using Apache Bench

```bash
# Load test the health endpoint
ab -n 1000 -c 10 http://localhost:4441/ontology_bot/v1/health/
```

### Using HTTPie

```bash
# Pretty-print response
http GET localhost:4441/ontology_bot/v1/health

# With timing
http --pretty=all --print=HhBb GET localhost:4441/ontology_bot/v1/health
```

## HTTP Status Codes

- **200 OK**: Service is healthy and responding normally
- **503 Service Unavailable**: Service is down or unable to respond (future enhancement)

## Future Enhancements

The readiness check will be enhanced to verify:
- Database connectivity (PostgreSQL, ClickHouse)
- Redis availability
- External service dependencies
- Cache status
