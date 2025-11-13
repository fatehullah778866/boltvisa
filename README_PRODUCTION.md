# 🚀 Production Deployment Guide

## Pre-Deployment Checklist

### Security
- [x] Rate limiting implemented
- [x] Password reset flow
- [x] Audit logging enabled
- [x] Input validation
- [x] Error boundaries
- [ ] HTTPS/TLS configured
- [ ] Security headers configured
- [ ] Penetration testing completed
- [ ] Dependency vulnerabilities scanned

### Performance
- [x] Load testing scripts ready
- [x] Metrics collection enabled
- [x] Health checks configured
- [ ] Load testing executed
- [ ] Performance baseline documented
- [ ] Database indexes optimized

### Monitoring
- [x] Logging middleware
- [x] Metrics endpoint
- [x] Audit logging
- [ ] Error tracking (Sentry) configured
- [ ] Alerting rules configured
- [ ] Dashboard created

### Infrastructure
- [x] Kubernetes configs ready
- [x] Cloud Build pipeline configured
- [x] Docker images built
- [ ] Secrets management configured
- [ ] Backup strategy implemented
- [ ] Disaster recovery plan documented

## Deployment Steps

### 1. Environment Setup

```bash
# Set environment variables
export GCP_PROJECT_ID=your-project-id
export DATABASE_URL=postgres://...
export JWT_SECRET=$(openssl rand -hex 32)
export SENDGRID_API_KEY=...
export STRIPE_SECRET_KEY=...
```

### 2. Database Setup

```bash
# Create Cloud SQL instance
gcloud sql instances create boltvisa-db \
  --database-version=POSTGRES_14 \
  --tier=db-f1-micro \
  --region=us-central1

# Create database
gcloud sql databases create boltvisa --instance=boltvisa-db

# Run migrations (automatic on startup)
```

### 3. Build and Push Images

```bash
# Build backend
cd apps/api
docker build -t gcr.io/$GCP_PROJECT_ID/boltvisa-api:latest .

# Build frontend
cd ../..
docker build -f apps/web/Dockerfile -t gcr.io/$GCP_PROJECT_ID/boltvisa-web:latest .

# Push images
docker push gcr.io/$GCP_PROJECT_ID/boltvisa-api:latest
docker push gcr.io/$GCP_PROJECT_ID/boltvisa-web:latest
```

### 4. Deploy to Cloud Run

```bash
# Deploy backend
gcloud run deploy boltvisa-api \
  --image gcr.io/$GCP_PROJECT_ID/boltvisa-api:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars DATABASE_URL=$DATABASE_URL,JWT_SECRET=$JWT_SECRET

# Deploy frontend
gcloud run deploy boltvisa-web \
  --image gcr.io/$GCP_PROJECT_ID/boltvisa-web:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars NEXT_PUBLIC_API_URL=https://boltvisa-api-xxx.run.app
```

### 5. Deploy to Kubernetes (Alternative)

```bash
# Update deployment.yaml with your project ID
sed -i 's/PROJECT_ID/your-project-id/g' k8s/deployment.yaml

# Create secrets
kubectl create secret generic boltvisa-secrets \
  --from-literal=database-url=$DATABASE_URL \
  --from-literal=jwt-secret=$JWT_SECRET

# Deploy
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/frontend-deployment.yaml
```

## Post-Deployment Verification

### 1. Health Checks

```bash
# Backend health
curl https://your-api-url/health

# Frontend
curl https://your-frontend-url/
```

### 2. API Testing

```bash
# Test registration
curl -X POST https://your-api-url/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","first_name":"Test","last_name":"User"}'

# Test login
curl -X POST https://your-api-url/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### 3. Load Testing

```bash
# Install k6
brew install k6  # macOS
# or download from https://k6.io/

# Run load test
k6 run load-test.js --env API_URL=https://your-api-url
```

### 4. Monitoring

```bash
# Check metrics
curl https://your-api-url/metrics

# View logs
gcloud logging read "resource.type=cloud_run_revision" --limit 50
```

## Security Hardening

### 1. Enable HTTPS

Cloud Run automatically provides HTTPS. Ensure:
- Custom domain configured (if using)
- SSL certificate valid
- Redirect HTTP to HTTPS

### 2. Security Headers

Configure at load balancer level:
```
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
```

### 3. Rate Limiting

Current implementation uses in-memory rate limiting. For production:
- Consider Redis-based rate limiting for distributed systems
- Configure Cloud Armor for DDoS protection
- Set up WAF rules

### 4. Secrets Management

Use Google Secret Manager:
```bash
# Create secrets
echo -n "your-secret" | gcloud secrets create jwt-secret --data-file=-

# Grant access
gcloud secrets add-iam-policy-binding jwt-secret \
  --member="serviceAccount:your-sa@project.iam.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"
```

## Monitoring & Alerts

### 1. Set Up Alerts

```bash
# High error rate
gcloud alpha monitoring policies create \
  --notification-channels=CHANNEL_ID \
  --display-name="High Error Rate" \
  --condition-display-name="Error rate > 5%" \
  --condition-threshold-value=0.05

# High latency
gcloud alpha monitoring policies create \
  --notification-channels=CHANNEL_ID \
  --display-name="High Latency" \
  --condition-display-name="p95 latency > 1s" \
  --condition-threshold-value=1000
```

### 2. Dashboard

Create Grafana dashboard or use Cloud Monitoring:
- Request rate
- Error rate
- Latency (p50, p95, p99)
- Database performance
- System resources

## Backup & Recovery

### Database Backups

```bash
# Enable automated backups
gcloud sql instances patch boltvisa-db \
  --backup-start-time=03:00 \
  --enable-bin-log

# Manual backup
gcloud sql backups create --instance=boltvisa-db
```

### Disaster Recovery

1. **Database:**
   - Automated daily backups
   - Point-in-time recovery
   - Cross-region replication (optional)

2. **Application:**
   - Multi-region deployment
   - Blue-green deployments
   - Rollback procedures

## Rollback Procedure

### Cloud Run

```bash
# List revisions
gcloud run revisions list --service=boltvisa-api

# Rollback to previous revision
gcloud run services update-traffic boltvisa-api \
  --to-revisions=REVISION_NAME=100
```

### Kubernetes

```bash
# Rollback deployment
kubectl rollout undo deployment/boltvisa-api

# Check rollout status
kubectl rollout status deployment/boltvisa-api
```

## Performance Optimization

### Database

```sql
-- Add indexes for frequently queried fields
CREATE INDEX idx_applications_user_id ON visa_applications(user_id);
CREATE INDEX idx_applications_status ON visa_applications(status);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
```

### Caching

Consider adding Redis for:
- Frequently accessed visa categories
- User session data
- Rate limiting (distributed)

### CDN

Configure Cloud CDN for:
- Static assets
- Images
- API responses (if appropriate)

## Troubleshooting

### Common Issues

1. **Database Connection Errors**
   - Check Cloud SQL instance status
   - Verify connection string
   - Check firewall rules

2. **Rate Limiting Issues**
   - Check rate limit configuration
   - Monitor `/metrics` endpoint
   - Adjust limits if needed

3. **Authentication Errors**
   - Verify JWT_SECRET matches
   - Check token expiration
   - Review audit logs

4. **Performance Issues**
   - Check database query performance
   - Review metrics endpoint
   - Analyze slow queries

## Support

For production issues:
1. Check logs: `gcloud logging read`
2. Review metrics: `/metrics` endpoint
3. Check audit logs: Admin dashboard
4. Review health checks: `/health` endpoint

## Next Steps

After successful deployment:
1. Monitor for 24-48 hours
2. Review metrics and logs
3. Gather user feedback
4. Plan optimizations
5. Schedule regular security audits

