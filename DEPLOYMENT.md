# Payment Service Deployment Guide

This guide covers deploying the payment service to your k3s cluster via ArgoCD.

## Prerequisites

- Access to k3s cluster (root@37.27.40.86)
- ArgoCD running at https://argocd.vibeoholic.com
- Stripe API keys (test or production)
- GitHub Personal Access Token for pulling images

## Step 1: Create Kubernetes Namespace

```bash
ssh root@37.27.40.86
kubectl create namespace payment-service
```

## Step 2: Create Image Pull Secret

```bash
kubectl create secret docker-registry ghcr-pull-secret \
  --docker-server=ghcr.io \
  --docker-username=frallan97 \
  --docker-password=<GITHUB_PAT> \
  --namespace=payment-service
```

## Step 3: Create Application Secrets

Create a file `payment-service-secrets.yaml` with your actual values:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: payment-service-secrets
  namespace: payment-service
type: Opaque
stringData:
  stripe-api-key: "sk_test_..." # or sk_live_... for production
  stripe-webhook-secret: "whsec_..."
  database-url: "postgresql://paymentuser:CHANGEME@payment-service-postgresql:5432/paymentdb?sslmode=disable"
  redis-url: "redis://payment-service-redis-master:6379"
```

Apply the secret:
```bash
kubectl apply -f payment-service-secrets.yaml
rm payment-service-secrets.yaml  # Don't commit this!
```

## Step 4: Add to ArgoCD ApplicationSet

On your local machine, edit the k3s-infra repository:

```bash
cd /path/to/k3s-infra
```

Edit `clusters/main/apps/app-of-apps.yaml` and add:

```yaml
  - name: payment-service
    namespace: payment-service
    repoURL: https://github.com/Frallan97/payment-service.git
    targetRevision: main
    path: charts/payment-service
    values: |
      image:
        repository: ghcr.io/frallan97/payment-service
        tag: latest

      ingress:
        enabled: true
        className: traefik
        annotations:
          cert-manager.io/cluster-issuer: letsencrypt-prod
        hosts:
          - host: payments.vibeoholic.com
            paths:
              - path: /
                pathType: Prefix
        tls:
          - secretName: payment-service-tls
            hosts:
              - payments.vibeoholic.com

      postgresql:
        enabled: true
        auth:
          username: paymentuser
          password: CHANGEME  # Use a strong password
          database: paymentdb

      redis:
        enabled: true

      resources:
        limits:
          cpu: 500m
          memory: 512Mi
        requests:
          cpu: 250m
          memory: 256Mi
```

Commit and push:
```bash
git add clusters/main/apps/app-of-apps.yaml
git commit -m "Add payment-service to ArgoCD"
git push origin main
```

## Step 5: Wait for ArgoCD Sync

ArgoCD will automatically detect the change and deploy:

1. Go to https://argocd.vibeoholic.com
2. Find the `payment-service` application
3. Click "Sync" if needed
4. Monitor the deployment progress

## Step 6: Verify Deployment

```bash
# Check pods
kubectl get pods -n payment-service

# Check services
kubectl get svc -n payment-service

# Check ingress
kubectl get ingress -n payment-service

# View logs
kubectl logs -n payment-service -l app.kubernetes.io/name=payment-service --tail=100

# Test health endpoint
curl https://payments.vibeoholic.com/health
```

Expected response:
```
OK
```

## Step 7: Test the API

### Get JWT Token from Auth Service

First, authenticate with your auth service to get a JWT token.

### Create a Test Payment

```bash
curl -X POST https://payments.vibeoholic.com/api/payments \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $(uuidgen)" \
  -d '{
    "provider": "stripe",
    "amount": 10000,
    "currency": "SEK",
    "description": "Test payment"
  }'
```

### Check Metrics

```bash
curl https://payments.vibeoholic.com/metrics
```

## Step 8: Configure Stripe Webhooks

1. Go to Stripe Dashboard → Developers → Webhooks
2. Add endpoint: `https://payments.vibeoholic.com/api/webhooks/stripe`
3. Select events to listen to:
   - `payment_intent.succeeded`
   - `payment_intent.payment_failed`
   - `payment_intent.canceled`
   - `customer.subscription.created`
   - `customer.subscription.updated`
   - `customer.subscription.deleted`
4. Copy the webhook signing secret
5. Update your Kubernetes secret:

```bash
kubectl edit secret payment-service-secrets -n payment-service
# Update stripe-webhook-secret value (base64 encoded)
```

## Step 9: Set Up Monitoring

### Prometheus ServiceMonitor

The Helm chart includes a ServiceMonitor. Verify it's created:

```bash
kubectl get servicemonitor -n payment-service
```

### View Metrics in Prometheus

1. Go to your Prometheus instance
2. Query: `payment_service_payments_total`
3. Check for metrics being scraped

### Create Grafana Dashboard

Import the dashboard from `grafana-dashboard.json` (see monitoring section below).

## Troubleshooting

### Pods Not Starting

```bash
# Check pod status
kubectl describe pod -n payment-service <pod-name>

# Check logs
kubectl logs -n payment-service <pod-name>

# Check events
kubectl get events -n payment-service --sort-by='.lastTimestamp'
```

### Database Connection Issues

```bash
# Check PostgreSQL pod
kubectl get pods -n payment-service -l app.kubernetes.io/name=postgresql

# Test connection from payment service pod
kubectl exec -it -n payment-service <payment-pod> -- /bin/sh
# Inside pod: try to connect to database
```

### Ingress Not Working

```bash
# Check ingress
kubectl describe ingress -n payment-service payment-service

# Check Traefik logs
kubectl logs -n kube-system -l app.kubernetes.io/name=traefik

# Verify DNS
dig payments.vibeoholic.com
```

### TLS Certificate Issues

```bash
# Check certificate
kubectl get certificate -n payment-service

# Check cert-manager logs
kubectl logs -n cert-manager -l app=cert-manager
```

## Rollback

If you need to rollback:

```bash
# Via ArgoCD UI
# 1. Go to payment-service app
# 2. Click "History and Rollback"
# 3. Select previous revision
# 4. Click "Rollback"

# Via kubectl
kubectl rollout undo deployment/payment-service -n payment-service
```

## Scaling

### Manual Scaling

```bash
kubectl scale deployment payment-service -n payment-service --replicas=3
```

### Enable Horizontal Pod Autoscaler

Edit `values.yaml` in ArgoCD:
```yaml
autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
```

## Monitoring Queries

### Useful Prometheus Queries

```promql
# Request rate
rate(payment_service_http_requests_total[5m])

# Error rate
rate(payment_service_http_requests_total{status=~"5.."}[5m])

# Payment success rate
rate(payment_service_payments_total{status="succeeded"}[5m])
  /
rate(payment_service_payments_total[5m])

# P95 latency
histogram_quantile(0.95, rate(payment_service_http_request_duration_seconds_bucket[5m]))

# Active subscriptions
payment_service_subscriptions_active
```

## Security Considerations

1. **Rotate Secrets Regularly**: Update Stripe API keys and database passwords periodically
2. **Use Network Policies**: Restrict pod-to-pod communication
3. **Enable Pod Security Standards**: Use restricted security contexts
4. **Audit Logs**: Review ArgoCD and Kubernetes audit logs
5. **HTTPS Only**: Ensure all traffic uses TLS

## Maintenance

### Update Application

1. Push code changes to GitHub
2. GitHub Actions builds and pushes new image
3. Update image tag in ArgoCD (or use `latest` for auto-update)
4. ArgoCD syncs the changes

### Database Migrations

Migrations run automatically on pod startup. To run manually:

```bash
kubectl exec -it -n payment-service <pod-name> -- /bin/sh
# Inside pod, run migrations
```

### Backup Database

```bash
kubectl exec -it -n payment-service payment-service-postgresql-0 -- pg_dump -U paymentuser paymentdb > backup.sql
```

## Production Checklist

- [ ] Use production Stripe keys
- [ ] Set strong database passwords
- [ ] Configure proper resource limits
- [ ] Enable autoscaling
- [ ] Set up monitoring alerts
- [ ] Configure log aggregation
- [ ] Enable network policies
- [ ] Set up regular backups
- [ ] Document runbooks
- [ ] Test disaster recovery

## Support

For issues:
1. Check logs: `kubectl logs -n payment-service -l app.kubernetes.io/name=payment-service`
2. Check ArgoCD: https://argocd.vibeoholic.com
3. Check Prometheus metrics: `/metrics` endpoint
4. Review GitHub Actions: https://github.com/Frallan97/payment-service/actions
