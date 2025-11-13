# Testing Guide

## Backend Testing

### Unit Tests

Run backend unit tests:

```bash
cd apps/api
go test -v ./...
```

Run with coverage:

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Structure

- Test files: `*_test.go`
- Test functions: `TestXxx(t *testing.T)`
- Uses `testify` for assertions
- Uses SQLite in-memory database for testing

### Example Test

```go
func TestRegister(t *testing.T) {
    db := setupTestDB()
    cfg := &config.Config{JWTSecret: "test-secret"}
    h := New(db, cfg)
    // ... test implementation
}
```

## Frontend Testing

### Linting

```bash
npm run lint --workspace=apps/web
```

### Type Checking

```bash
npm run type-check --workspace=apps/web
```

## E2E Testing

### Setup

Install Playwright:

```bash
npx playwright install
```

### Run Tests

```bash
# Run all E2E tests
npm run test:e2e

# Run with UI
npm run test:e2e:ui

# Run specific test file
npx playwright test e2e/auth.spec.ts
```

### Test Files

- `e2e/auth.spec.ts` - Authentication flow tests
- `e2e/applications.spec.ts` - Application management tests

### Writing E2E Tests

```typescript
import { test, expect } from '@playwright/test';

test('should do something', async ({ page }) => {
  await page.goto('/path');
  await expect(page.locator('selector')).toBeVisible();
});
```

## CI/CD Testing

Tests run automatically on:
- Push to `main` or `develop`
- Pull requests

Test results are uploaded to:
- Codecov for coverage
- GitHub Actions artifacts for E2E reports

## Load Testing

### Using k6

```bash
# Install k6
brew install k6  # macOS
# or download from https://k6.io/

# Run load test
k6 run load-test.js
```

### Example Load Test Script

```javascript
import http from 'k6/http';
import { check } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 20 },
    { duration: '1m', target: 50 },
    { duration: '30s', target: 0 },
  ],
};

export default function () {
  const res = http.get('http://localhost:8080/health');
  check(res, { 'status is 200': (r) => r.status === 200 });
}
```

## Performance Testing

### Metrics Endpoint

Access metrics at `/metrics`:

```bash
curl http://localhost:8080/metrics
```

Returns:
```json
{
  "metrics": {
    "/api/v1/applications": {
      "count": 1234,
      "latency": 45,
      "errors": 2
    }
  }
}
```

## Best Practices

1. **Unit Tests**: Test individual functions/handlers in isolation
2. **Integration Tests**: Test API endpoints with database
3. **E2E Tests**: Test complete user flows
4. **Load Tests**: Test system under load
5. **Coverage**: Aim for >80% code coverage

## Debugging Tests

### Backend

```bash
# Run with verbose output
go test -v ./...

# Run specific test
go test -v -run TestRegister
```

### Frontend E2E

```bash
# Run in headed mode
npx playwright test --headed

# Debug mode
npx playwright test --debug
```

