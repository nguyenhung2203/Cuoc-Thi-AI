// k6 load/stability test for the Python AI service (/api/v1/generate).
//
// Run the AI service in mock mode to test throughput without spending tokens:
//   AI_MOCK=true uvicorn app.main:app --port 8000
//   k6 run -e AI=http://localhost:8000 tests/load/k6_ai_generate.js

import http from 'k6/http';
import { check, sleep } from 'k6';

const AI = __ENV.AI || 'http://localhost:8000';
const VUS = parseInt(__ENV.VUS || '20', 10);
const DURATION = __ENV.DURATION || '30s';

export const options = {
  scenarios: {
    steady: {
      executor: 'constant-vus',
      vus: VUS,
      duration: DURATION,
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<1500'],
    http_req_failed: ['rate<0.02'],
  },
};

export default function () {
  const health = http.get(`${AI}/health`);
  check(health, { 'ai health 200': (r) => r.status === 200 });

  const body = JSON.stringify({
    prompt: 'Sinh danh sách câu hỏi phỏng vấn (questions) cho vị trí backend developer.',
    model: 'gemini-1.5-flash',
    temperature: 0.2,
    max_tokens: 1024,
  });
  const res = http.post(`${AI}/api/v1/generate`, body, {
    headers: { 'Content-Type': 'application/json' },
  });
  check(res, {
    'generate 200': (r) => r.status === 200,
    'has data envelope': (r) => r.json('data') !== undefined,
  });
  sleep(0.5);
}
