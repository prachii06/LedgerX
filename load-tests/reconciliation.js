import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = __ENV.API_URL || 'http://localhost:8080';

export const options = {
    stages: [
        { duration: '30s', target: 10 }, 
        { duration: '1m', target: 30 },
        { duration: '30s', target: 0 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<800'],
        http_req_failed: ['rate<0.05'],
    },
};

export default function () {
    const payload = JSON.stringify({
        external_id: `rec-${__VU}-${__ITER}-${Date.now()}`,
        amount: 250.0,
        currency: 'EUR',
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    // 1. Create transaction
    const res = http.post(`${BASE_URL}/transactions`, payload, params);

    if (check(res, { 'transaction created': (r) => r.status === 201 })) {
        const txId = res.json('id');
        
        // 2. Wait for asynchronous processing (Kafka -> PG -> Reconcile)
        sleep(2);

        // 3. Check reconciliation status
        const recRes = http.get(`${BASE_URL}/reconcile/${txId}`);
        check(recRes, {
            'reconciliation fetched': (r) => r.status === 200,
            'status present': (r) => r.json('status') !== undefined,
        });
    }

    sleep(1);
}
