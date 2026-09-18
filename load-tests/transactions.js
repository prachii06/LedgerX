import http from 'k6/http';
import { check, sleep } from 'k6';

// Read target URL from env or use default
const BASE_URL = __ENV.API_URL || 'http://localhost:8080';

export const options = {
    stages: [
        { duration: '30s', target: 10 }, // Baseline
        { duration: '1m', target: 50 },  // Moderate
        { duration: '1m', target: 100 }, // High
        { duration: '30s', target: 0 },  // Ramp-down
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
        http_req_failed: ['rate<0.01'],   // Error rate should be less than 1%
    },
    summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(90)', 'p(95)', 'p(99)'],
};

export default function () {
    const payload = JSON.stringify({
        external_id: `ext-${__VU}-${__ITER}-${Date.now()}`,
        amount: Math.round((Math.random() * 1000 + 10) * 100) / 100,
        currency: 'USD',
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const res = http.post(`${BASE_URL}/transactions`, payload, params);

    check(res, {
        'status is 201': (r) => r.status === 201,
        'has id': (r) => r.json('id') !== undefined,
    });

    sleep(1);
}
