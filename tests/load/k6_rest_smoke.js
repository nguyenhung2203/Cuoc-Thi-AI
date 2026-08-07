// k6 REST load/stability test for the AI Interview Platform backend.
//
// Usage:
//   k6 run tests/load/k6_rest_smoke.js
//   k6 run -e API=http://localhost:8080 -e VUS=50 -e DURATION=1m tests/load/k6_rest_smoke.js
//
// Exercises the public + auth flow: health, register, login, and an
// authenticated read. Thresholds enforce latency and error-rate SLOs so the
// run FAILS (non-zero exit) if the backend is unstable under load.

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const API = __ENV.API || 'http://localhost:8080';
const BASE = `${API}/api/v1`;
const VUS = parseInt(__ENV.VUS || '20', 10);
const DURATION = __ENV.DURATION || '30s';

const errorRate = new Rate('business_errors');

export const options = {
  scenarios: {
    ramp: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '10s', target: VUS },
        { duration: DURATION, target: VUS },
        { duration: '10s', target: 0 },
      ],
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<800'], // 95% of requests under 800ms
    http_req_failed: ['rate<0.01'],   // <1% transport failures
    business_errors: ['rate<0.05'],   // <5% business-level errors
  },
};

function uniqueEmail() {
  return `load_${__VU}_${__ITER}_${Date.now()}@example.com`;
}

export default function () {
  // 1. Health check (must always be fast + healthy).
  const health = http.get(`${API}/health`);
  check(health, { 'health 200': (r) => r.status === 200 }) || errorRate.add(1);

  // 2. Register a fresh recruiter.
  const email = uniqueEmail();
  const payload = JSON.stringify({
    email,
    password: 'Passw0rd!23',
    full_name: 'Load Test User',
    role: 'recruiter',
  });
  const headers = { 'Content-Type': 'application/json' };

  const reg = http.post(`${BASE}/auth/register`, payload, { headers });
  const registered = check(reg, {
    'register 200/201': (r) => r.status === 200 || r.status === 201,
  });
  if (!registered) {
    errorRate.add(1);
    sleep(1);
    return;
  }

  // 3. Login and read /auth/me.
  const login = http.post(`${BASE}/auth/login`, JSON.stringify({ email, password: 'Passw0rd!23' }), { headers });
  const token = login.json('data.access_token') || login.json('access_token');
  check(login, { 'login 200': (r) => r.status === 200 }) || errorRate.add(1);

  if (token) {
    const me = http.get(`${BASE}/auth/me`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    check(me, { 'me 200': (r) => r.status === 200 }) || errorRate.add(1);
  }

  sleep(1);
}
