# Final Deployment Steps

The payment service is now ready for deployment. Follow these steps to deploy to your k3s cluster.

## Step 1: Add to ArgoCD ApplicationSet

Navigate to your k3s-infra repository and edit the ApplicationSet configuration:

```bash
cd /path/to/k3s-infra
```

Edit `clusters/main/apps/app-of-apps.yaml` and add this entry to the applications list:

```yaml
- name: payment-service
  namespace: payment-service
  repoURL: https://github.com/Frallan97/payment-service.git
  targetRevision: main
  path: charts/payment-service
  values: |
    replicaCount: 2

    image:
      repository: ghcr.io/frallan97/payment-service
      tag: latest
      pullPolicy: Always

    imagePullSecrets:
      - name: ghcr-pull-secret

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

    env:
      - name: PORT
        value: "8080"
      - name: ENV
        value: "production"
      - name: AUTH_SERVICE_URL
        value: "https://auth.vibeoholic.com"
      - name: ALLOWED_ORIGINS
        value: "https://app.vibeoholic.com,https://vibeoholic.com"

    secrets:
      - name: STRIPE_API_KEY
        valueFrom:
          secretKeyRef:
            name: payment-service-secrets
            key: stripe-api-key
      - name: STRIPE_WEBHOOK_SECRET
        valueFrom:
          secretKeyRef:
            name: payment-service-secrets
            key: stripe-webhook-secret
      - name: DATABASE_URL
        valueFrom:
          secretKeyRef:
            name: payment-service-secrets
            key: database-url
      - name: REDIS_URL
        valueFrom:
          secretKeyRef:
            name: payment-service-secrets
            key: redis-url

    postgresql:
      enabled: true
      auth:
        username: paymentuser
        password: CHANGE_THIS_PASSWORD
        database: paymentdb
      primary:
        persistence:
          enabled: true
          size: 10Gi

    redis:
      enabled: true
      auth:
        enabled: false
      master:
        persistence:
          enabled: true
          size: 2Gi

    resources:
      limits:
        cpu: 500m
        memory: 512Mi
      requests:
        cpu: 250m
        memory: 256Mi

    serviceMonitor:
      enabled: true
      interval: 30s
```

**Important:** Generate a strong random password for the PostgreSQL database and replace `CHANGE_THIS_PASSWORD` above.

## Step 2: Commit and Push to k3s-infra

```bash
cd /path/to/k3s-infra
git add clusters/main/apps/app-of-apps.yaml
git commit -m "Add payment-service to ArgoCD"
git push origin main
```

ArgoCD will automatically detect the change and begin deployment.

## Step 3: Create Kubernetes Secrets

Use the interactive deployment script to create the required secrets:

```bash
cd /home/frans-sjostrom/Documents/hezner-hosted-projects/payment-service
./scripts/deploy.sh
```

Select option 3 (Create application secrets) and provide:
- **Stripe API Key**: Your production Stripe secret key (sk_live_...) or test key (sk_test_...)
- **Stripe Webhook Secret**: Will be obtained in Step 5 after deployment
- **Database Password**: Same password you used in Step 1

## Step 4: Create Image Pull Secret

Run the deployment script again and select option 2 (Create image pull secret):

```bash
./scripts/deploy.sh
```

Provide:
- **GitHub Username**: frallan97
- **GitHub PAT**: Your GitHub Personal Access Token with read:packages permission

## Step 5: Monitor Deployment

Check deployment status using the script:

```bash
./scripts/deploy.sh
```

Select option 4 (Check deployment status) to see:
- Pod status
- Service endpoints
- Ingress configuration
- Recent events

Wait until all pods are Running and Ready (2/2).

## Step 6: Verify Service is Accessible

Test the health endpoint:

```bash
curl https://payments.vibeoholic.com/health
# Expected response: OK

curl https://payments.vibeoholic.com/metrics | head -20
# Expected: Prometheus metrics output
```

## Step 7: Configure Stripe Webhooks

1. Go to [Stripe Dashboard → Webhooks](https://dashboard.stripe.com/webhooks)
2. Click "Add endpoint"
3. Enter URL: `https://payments.vibeoholic.com/api/webhooks/stripe`
4. Select events to listen for:
   - `payment_intent.succeeded`
   - `payment_intent.payment_failed`
   - `payment_intent.canceled`
   - `customer.subscription.created`
   - `customer.subscription.updated`
   - `customer.subscription.deleted`
5. Click "Add endpoint"
6. Copy the **Signing secret** (starts with `whsec_`)
7. Update the Kubernetes secret:

```bash
ssh root@37.27.40.86
kubectl patch secret payment-service-secrets \
  -n payment-service \
  --type='json' \
  -p='[{"op": "replace", "path": "/data/stripe-webhook-secret", "value":"'$(echo -n "whsec_YOUR_SECRET" | base64)'"}]'

# Restart pods to pick up the new secret
kubectl rollout restart deployment payment-service -n payment-service
```

## Step 8: Import Grafana Dashboard

1. Open your Grafana instance
2. Navigate to Dashboards → Import
3. Upload the file: `monitoring/grafana-dashboard.json`
4. Select your Prometheus data source
5. Click Import

The dashboard includes:
- Payment rate and success rate
- Active subscriptions count
- HTTP request rate by endpoint
- Request latency (P50, P95)
- Payments by provider and status
- Webhook events and errors
- HTTP 5xx error rate

## Step 9: Test the Service

Create a test payment:

```bash
# First, get a JWT token from your auth-service
# Then create a payment:

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

Expected response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "provider": "stripe",
  "provider_payment_id": "pi_3ABC123",
  "amount": 10000,
  "currency": "SEK",
  "status": "pending",
  "client_secret": "pi_3ABC123_secret_xyz",
  "created_at": "2026-01-12T10:00:00Z"
}
```

## Step 10: Configure DNS (if needed)

If `payments.vibeoholic.com` is not yet configured, add a DNS A record:

```
Type: A
Name: payments
Value: 37.27.40.86 (your k3s master node IP)
Proxied: No (to allow cert-manager to issue certificates)
```

Wait for DNS propagation (usually 5-15 minutes).

## Troubleshooting

### Pods Not Starting

```bash
# View pod details
ssh root@37.27.40.86
kubectl describe pod -n payment-service -l app.kubernetes.io/name=payment-service

# Check logs
kubectl logs -n payment-service -l app.kubernetes.io/name=payment-service --tail=100
```

Common issues:
- Image pull errors → Check ghcr-pull-secret is created
- Database connection errors → Check PostgreSQL pod is running
- Secret errors → Verify payment-service-secrets exists with correct keys

### Ingress 404 Errors

```bash
# Check ingress configuration
kubectl get ingress -n payment-service payment-service -o yaml

# Verify TLS certificate
kubectl get certificate -n payment-service payment-service-tls
```

### Database Connection Issues

```bash
# Check PostgreSQL pod
kubectl get pods -n payment-service -l app.kubernetes.io/name=postgresql

# Check PostgreSQL logs
kubectl logs -n payment-service -l app.kubernetes.io/name=postgresql
```

## Post-Deployment Checklist

- [ ] All pods running (2 payment-service replicas + PostgreSQL + Redis)
- [ ] Health endpoint responding (https://payments.vibeoholic.com/health)
- [ ] Metrics endpoint available (https://payments.vibeoholic.com/metrics)
- [ ] TLS certificate issued (no browser warnings)
- [ ] Stripe webhooks configured and receiving events
- [ ] Grafana dashboard imported and displaying data
- [ ] Test payment created successfully
- [ ] Test subscription created successfully
- [ ] Webhook events processing correctly
- [ ] Rate limiting working (check X-RateLimit headers)
- [ ] ArgoCD shows sync status as "Healthy" and "Synced"

## Monitoring URLs

- **Service**: https://payments.vibeoholic.com
- **Health Check**: https://payments.vibeoholic.com/health
- **Metrics**: https://payments.vibeoholic.com/metrics
- **ArgoCD**: https://argocd.vibeoholic.com/applications/payment-service
- **Stripe Dashboard**: https://dashboard.stripe.com

## Next Steps (Optional Enhancements)

1. **Add Swish Provider**: Implement Swish integration with TLS certificates
2. **Add More Tests**: Integration tests, E2E tests with testcontainers
3. **Add OpenAPI Documentation**: Swagger UI for API exploration
4. **Add Alerting**: Configure Prometheus alerts for critical metrics
5. **Add Payment Analytics**: Dashboard for business metrics
6. **Add Webhook Retry Logic**: Exponential backoff for failed webhook processing
7. **Add Payment Dispute Handling**: Support for Stripe dispute events

---

**Deployment Status**: Ready ✅

All code is complete, tested, and ready for production deployment!
