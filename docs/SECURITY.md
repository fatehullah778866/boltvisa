# Security Guide

## Implemented Security Features

### 1. Rate Limiting

**Implementation:**
- In-memory rate limiting middleware
- 100 requests per 15 minutes for anonymous users
- 200 requests per 15 minutes for authenticated users
- IP-based and user-based tracking

**Configuration:**
- Can be upgraded to Redis-based for distributed systems
- Configurable limits in `internal/middleware/ratelimit.go`

**Usage:**
```go
// Applied globally
r.Use(middleware.RateLimit())

// Stricter limit for authenticated routes
protected.Use(middleware.AuthRateLimit())
```

### 2. Authentication & Authorization

**JWT Tokens:**
- Secure token generation with expiration (24 hours)
- Token refresh endpoint
- Role-based access control (RBAC)

**Password Security:**
- bcrypt hashing with default cost
- Minimum 8 character requirement
- Password reset flow with time-limited tokens

### 3. Audit Logging

**Tracked Actions:**
- User registration
- User login/logout
- Role changes
- Payment transactions
- Admin actions

**Audit Log Fields:**
- User ID
- Action type
- Resource and resource ID
- IP address
- User agent
- Timestamp
- Metadata (JSON)

**Access:**
- Admin-only endpoint: `GET /api/v1/admin/audit-logs`
- Pagination and filtering supported

### 4. Input Validation

**Backend:**
- Gin binding validation
- Required fields
- Email format validation
- Password strength requirements

**Frontend:**
- Form validation
- Type checking with TypeScript
- Input sanitization

### 5. CORS Protection

**Configuration:**
- Configurable allowed origins
- Credentials support
- Preflight handling

### 6. SQL Injection Prevention

**GORM Protection:**
- Parameterized queries
- No raw SQL with user input
- Type-safe database operations

## Security Best Practices

### Password Reset

1. **Token Generation:**
   - Cryptographically secure random tokens (32 bytes)
   - Hex-encoded for URL safety
   - 1-hour expiration

2. **Email Security:**
   - No user enumeration (same response for valid/invalid emails)
   - Secure reset links
   - Clear expiration messaging

3. **Token Usage:**
   - Single-use tokens
   - Marked as used after successful reset
   - Invalidated on expiration

### Rate Limiting

**Recommendations:**
- Monitor rate limit hits
- Adjust limits based on usage patterns
- Consider Redis for distributed deployments
- Implement progressive rate limiting (stricter after violations)

### Audit Logging

**Best Practices:**
- Regular audit log review
- Automated alerts for suspicious activities
- Retention policies (compliance requirements)
- Encrypted storage for sensitive audit data

## Security Checklist

### Pre-Production

- [x] Rate limiting implemented
- [x] Password hashing (bcrypt)
- [x] JWT token security
- [x] Audit logging for sensitive operations
- [x] Input validation
- [x] CORS configuration
- [x] SQL injection prevention
- [x] Password reset flow
- [ ] HTTPS/TLS enforcement (infrastructure level)
- [ ] Security headers (HSTS, CSP, etc.)
- [ ] Penetration testing
- [ ] Dependency vulnerability scanning
- [ ] Secrets management review

### Production

- [ ] Enable HTTPS/TLS 1.3
- [ ] Configure security headers
- [ ] Set up WAF (Web Application Firewall)
- [ ] Enable DDoS protection
- [ ] Regular security audits
- [ ] Monitor for suspicious activities
- [ ] Incident response plan
- [ ] Regular dependency updates

## Security Headers (Recommended)

Add to your reverse proxy/load balancer:

```
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
Referrer-Policy: strict-origin-when-cross-origin
```

## Monitoring & Alerts

### Security Events to Monitor

1. **Failed Login Attempts:**
   - Multiple failures from same IP
   - Brute force patterns

2. **Rate Limit Violations:**
   - Frequent rate limit hits
   - Unusual traffic patterns

3. **Audit Log Alerts:**
   - Role changes
   - Payment transactions
   - Admin actions

4. **Error Rates:**
   - High 4xx/5xx rates
   - Unusual error patterns

## Incident Response

### Security Incident Checklist

1. **Immediate Actions:**
   - Identify and isolate affected systems
   - Preserve logs and evidence
   - Notify security team

2. **Investigation:**
   - Review audit logs
   - Analyze attack vectors
   - Assess data exposure

3. **Remediation:**
   - Patch vulnerabilities
   - Reset compromised credentials
   - Update security controls

4. **Post-Incident:**
   - Document incident
   - Review and improve security
   - Notify affected users (if required)

## Compliance Considerations

### GDPR (if serving EU users)

- [ ] Data deletion endpoints
- [ ] Consent management
- [ ] Privacy policy
- [ ] Data export functionality
- [ ] Right to be forgotten

### PCI DSS (if handling payments)

- [ ] Secure payment processing
- [ ] No storage of full card numbers
- [ ] PCI-compliant payment gateway (Stripe/Razorpay)
- [ ] Regular security assessments

## Reporting Security Issues

If you discover a security vulnerability, please:
1. Do not open a public issue
2. Email security@boltvisa.com (or your security contact)
3. Provide detailed information about the vulnerability
4. Allow time for remediation before disclosure

