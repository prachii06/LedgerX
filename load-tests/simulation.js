import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = __ENV.API_URL || 'http://localhost:8080';

export const options = {
    stages: [
        { duration: '30s', target: 5 }, // Baseline
        { duration: '1m', target: 10 }, // Moderate (lower because simulation generates multiple tx)
        { duration: '30s', target: 0 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<1000'],
        http_req_failed: ['rate<0.01'],
    },
};

export default function () {
    const payload = JSON.stringify({
        count: 5,
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const res = http.post(`${BASE_URL}/simulate`, payload, params);

    check(res, {
        'status is 200': (r) => r.status === 200,
    });

    sleep(2);
}
