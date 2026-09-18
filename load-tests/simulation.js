import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';

const BASE_URL = __ENV.API_URL || 'http://localhost:8080';

const simulatedTransactions = new Counter('simulated_transactions');

export const options = {
    stages: [
        { duration: '30s', target: 5 },
        { duration: '1m', target: 10 },
        { duration: '30s', target: 0 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<1000'],
        http_req_failed: ['rate<0.01'],
    },
    summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(90)', 'p(95)', 'p(99)'],
};

export default function () {
    const count = 5;
    const payload = JSON.stringify({
        count: count,
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const res = http.post(`${BASE_URL}/simulate`, payload, params);

    const isSuccessful = check(res, {
        'status is 200': (r) => r.status === 200,
        'transaction_ids exists': (r) => r.json('transaction_ids') !== undefined,
        'transaction_ids is array': (r) => {
            try {
                return Array.isArray(r.json('transaction_ids'));
            } catch (e) {
                return false;
            }
        },
        'transaction_ids.length === count': (r) => {
            try {
                const ids = r.json('transaction_ids');
                return Array.isArray(ids) && ids.length === count;
            } catch (e) {
                return false;
            }
        },
    });

    if (isSuccessful) {
        simulatedTransactions.add(res.json('count'));
    }

    sleep(2);
}
