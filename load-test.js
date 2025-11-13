import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const apiDuration = new Trend('api_duration');

export const options = {
  stages: [
    { duration: '30s', target: 20 },   // Ramp up to 20 users
    { duration: '1m', target: 50 },    // Stay at 50 users
    { duration: '30s', target: 100 },  // Ramp up to 100 users
    { duration: '1m', target: 100 },   // Stay at 100 users
    { duration: '30s', target: 0 },    // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'], // 95% of requests < 500ms, 99% < 1s
    http_req_failed: ['rate<0.01'],                 // Error rate < 1%
    errors: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.API_URL || 'http://localhost:8080';

// Test user credentials (should be created beforehand)
const TEST_EMAIL = __ENV.TEST_EMAIL || 'test@example.com';
const TEST_PASSWORD = __ENV.TEST_PASSWORD || 'password123';

let authToken = '';

export function setup() {
  // Login and get token
  const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`, JSON.stringify({
    email: TEST_EMAIL,
    password: TEST_PASSWORD,
  }), {
    headers: { 'Content-Type': 'application/json' },
  });

  if (loginRes.status === 200) {
    const body = JSON.parse(loginRes.body);
    return { token: body.token };
  }
  
  return { token: null };
}

export default function (data) {
  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${data.token}`,
  };

  // Test health endpoint
  let res = http.get(`${BASE_URL}/health`);
  check(res, {
    'health check status is 200': (r) => r.status === 200,
  });
  errorRate.add(res.status >= 400);
  apiDuration.add(res.timings.duration);
  sleep(1);

  // Test get current user
  res = http.get(`${BASE_URL}/api/v1/users/me`, { headers });
  check(res, {
    'get user status is 200': (r) => r.status === 200,
    'get user has user data': (r) => {
      const body = JSON.parse(r.body);
      return body.id !== undefined;
    },
  });
  errorRate.add(res.status >= 400);
  apiDuration.add(res.timings.duration);
  sleep(1);

  // Test get visa categories
  res = http.get(`${BASE_URL}/api/v1/visa-categories`, { headers });
  check(res, {
    'get categories status is 200': (r) => r.status === 200,
    'get categories returns array': (r) => {
      const body = JSON.parse(r.body);
      return Array.isArray(body);
    },
  });
  errorRate.add(res.status >= 400);
  apiDuration.add(res.timings.duration);
  sleep(1);

  // Test get applications
  res = http.get(`${BASE_URL}/api/v1/applications?page=1&page_size=10`, { headers });
  check(res, {
    'get applications status is 200': (r) => r.status === 200,
  });
  errorRate.add(res.status >= 400);
  apiDuration.add(res.timings.duration);
  sleep(1);

  // Test get notifications
  res = http.get(`${BASE_URL}/api/v1/notifications`, { headers });
  check(res, {
    'get notifications status is 200': (r) => r.status === 200,
  });
  errorRate.add(res.status >= 400);
  apiDuration.add(res.timings.duration);
  sleep(1);
}

export function teardown(data) {
  // Cleanup if needed
}

