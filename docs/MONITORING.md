# Monitoring & Observability Guide

## Metrics

### Application Metrics

Access metrics endpoint:

```bash
curl http://localhost:8080/metrics
```

Metrics include:
- Request count per endpoint
- Average latency per endpoint
- Error count per endpoint

### Custom Metrics

Add custom metrics in handlers:

```go
// Record custom metric
metricsMutex.Lock()
customMetric++
metricsMutex.Unlock()
```

## Logging

### Structured Logging

All requests are logged with:
- Timestamp
- HTTP method
- Path
- Status code
- Latency
- Client IP
- User agent

### Log Levels

- **INFO**: Normal operations
- **WARN**: Warning conditions
- **ERROR**: Error conditions

### Cloud Logging (Production)

In production, logs are automatically sent to:
- **Google Cloud Logging** (when deployed to GCP)
- **Stackdriver** (legacy name)

View logs:

```bash
gcloud logging read "resource.type=cloud_run_revision" --limit 50
```

## Health Checks

### Health Endpoint

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "ok"
}
```

### Kubernetes Health Checks

Configured in `k8s/deployment.yaml`:
- **Liveness Probe**: Checks if container is alive
- **Readiness Probe**: Checks if container is ready to serve traffic

## Error Tracking

### Error Handling

Errors are logged with:
- Error message
- Stack trace (in development)
- Request context
- User ID (if authenticated)

### Error Monitoring

In production, integrate with:
- **Sentry** (recommended)
- **Google Error Reporting**
- **Datadog**

## Performance Monitoring

### APM Tools

Recommended tools:
- **Google Cloud Trace** (for GCP)
- **New Relic**
- **Datadog APM**

### Database Monitoring

Monitor PostgreSQL:
- Query performance
- Connection pool usage
- Slow queries

### Cache Monitoring

Monitor Redis (if implemented):
- Hit/miss ratio
- Memory usage
- Eviction rate

## Alerting

### Alert Rules

Set up alerts for:
- High error rate (>5%)
- High latency (>1s p95)
- Low availability (<99%)
- Database connection issues

### Alert Channels

Configure alerts to:
- Email
- Slack
- PagerDuty
- SMS (via Twilio)

## Dashboards

### Grafana Dashboard

Create dashboards for:
- Request rate
- Error rate
- Latency (p50, p95, p99)
- Database performance
- System resources

### Google Cloud Console

View metrics in:
- Cloud Monitoring
- Cloud Logging
- Cloud Trace

## Best Practices

1. **Log Aggregation**: Centralize logs
2. **Structured Logging**: Use JSON format
3. **Correlation IDs**: Track requests across services
4. **Sampling**: Sample high-volume logs
5. **Retention**: Set appropriate retention policies

## Example Monitoring Stack

```
Application → Cloud Logging → BigQuery
           → Cloud Monitoring → Grafana
           → Cloud Trace → Performance Dashboard
```

