import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Trend } from 'k6/metrics';

const BASE_URL = __ENV.API_URL || 'http://localhost:8080';

const recMatched = new Counter('reconciliation_matched');
const recMissing = new Counter('reconciliation_missing');
const recDuplicate = new Counter('reconciliation_duplicate');
const recMismatch = new Counter('reconciliation_mismatch');
const recOutOfOrder = new Counter('reconciliation_out_of_order');
const recCurrencyMismatch = new Counter('reconciliation_currency_mismatch');
const recTimeout = new Counter('reconciliation_timeout');
const recTime = new Trend('reconciliation_time');

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
    summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(90)', 'p(95)', 'p(99)'],
};

export default function () {
    const payload = JSON.stringify({
        count: 1,
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const startTime = Date.now();
    const res = http.post(`${BASE_URL}/simulate`, payload, params);

    const simulationSuccessful = check(res, {
        'simulation HTTP request successful': (r) => r.status === 200,
        'transaction_ids exists': (r) => {
            try { return r.json('transaction_ids') !== undefined; }
            catch (e) { return false; }
        },
        'transaction_ids contains at least one ID': (r) => {
            try {
                const ids = r.json('transaction_ids');
                return Array.isArray(ids) && ids.length > 0;
            } catch (e) {
                return false;
            }
        },
    });

    if (simulationSuccessful) {
        const txId = res.json('transaction_ids')[0];
        
        let found = false;
        let attempts = 0;
        const maxAttempts = 10;
        
        while (attempts < maxAttempts && !found) {
            sleep(1); // poll interval
            attempts++;
            
            // Tagging the request helps keep the k6 summary clean
            const recRes = http.get(`${BASE_URL}/reconcile/${txId}`, {
                tags: { name: 'PollReconciliation' }
            });
            
            if (recRes.status === 200) {
                let status;
                try {
                    status = recRes.json('status');
                } catch (e) {
                    status = undefined;
                }

                if (status) {
                    found = true;
                    
                    check(recRes, {
                        'reconciliation HTTP request successful': (r) => r.status === 200,
                        'reconciliation result eventually available': (r) => true,
                    });
                    
                    const timeTaken = Date.now() - startTime;
                    recTime.add(timeTaken);
                    
                    switch (status) {
                        case 'MATCHED': recMatched.add(1); break;
                        case 'MISSING': recMissing.add(1); break;
                        case 'DUPLICATE': recDuplicate.add(1); break;
                        case 'MISMATCH': recMismatch.add(1); break;
                        case 'OUT_OF_ORDER': recOutOfOrder.add(1); break;
                        case 'CURRENCY_MISMATCH': recCurrencyMismatch.add(1); break;
                    }
                }
            }
        }
        
        if (!found) {
            check(null, {
                'reconciliation result eventually available': () => false,
            });
            recTimeout.add(1);
        }
    }
}
