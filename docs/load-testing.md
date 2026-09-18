# LedgerX Load Testing Report

## 1. Environment and Configuration
- **Test Environment:** Local machine (Windows).
- **k6 Version:** v2.2.0 (windows/amd64).
- **Database Configuration:** PostgreSQL running in Docker (`postgres` container, port 5433).
- **Kafka Configuration:** Kafka broker running in Docker.
- **Redis Configuration:** Redis running in Docker.

## 2. Test Objectives and Results

We use three completely decoupled tests to measure the performance of our APIs and event pipelines. 
*Note: `p99` latencies are explicitly measured in the scripts.*

### Test A: Transaction API (`transactions.js`)
- **Target:** `POST /transactions`
- **Goal:** Simulates typical external API client traffic, measuring the performance and reliability of the direct transaction creation API.
- **Events Expected:** **0** (This endpoint intentionally creates transactions without generating Kafka events).
- **Results:** 
  - Tracks success rate, HTTP errors, average latency, p90, p95, p99, max latency, and throughput.

### Test B: Event Simulation Pipeline (`simulation.js`)
- **Target:** `POST /simulate`
- **Goal:** Stresses the actual LedgerX distributed event and reconciliation pipeline by injecting delayed, missing, and duplicate events.
- **Metrics Tracked in k6:**
  - HTTP metrics: successful requests, average latency, p95, throughput.
  - Custom metrics: total simulated transactions (`successful requests * count`).
- **Interpretation:** This test seeds the system. Verification of Kafka throughput and Database persistency must be done post-test using the SQL queries provided below.

### Test C: Reconciliation Pipeline (`reconciliation.js`)
- **Target:** `POST /simulate` followed by `GET /reconcile/:id`
- **Goal:** Simulates full processing loop and verifies eventual reconciliation status polling.
- **Interpretation:** `MISSING` and `DUPLICATE` statuses are **expected** outcomes from the simulation pipeline. The test considers a status check successful as long as the API is responsive and returns a status, preventing false-positive test failures.

## 3. Database Verification (Post-Simulation SQL Queries)

After running Test B (`simulation.js`), use the following queries to verify ONLY the transactions created by the simulation test. This isolation avoids polluting metrics with `POST /transactions` data. 
*(Simulated transactions have their `external_id` prefixed with `SIM-%`).*

```sql
-- 1. Number of simulated transactions
SELECT COUNT(*) FROM transactions WHERE external_id LIKE 'SIM-%';

-- 2. Number of events produced for simulated transactions
SELECT COUNT(*) FROM transaction_events 
WHERE transaction_id IN (SELECT id FROM transactions WHERE external_id LIKE 'SIM-%');

-- 3. Number of unique simulated transactions in reconciliation_results
SELECT COUNT(*) FROM reconciliation_results 
WHERE transaction_id IN (SELECT id FROM transactions WHERE external_id LIKE 'SIM-%');

-- 4. Reconciliation Breakdown for simulated transactions
SELECT status, COUNT(*) 
FROM reconciliation_results 
WHERE transaction_id IN (SELECT id FROM transactions WHERE external_id LIKE 'SIM-%')
GROUP BY status ORDER BY status;
```

## 4. Reliability Observations and Conclusions

- **Correct Conclusions:** Previously, thousands of events were reported as "dropped". This conclusion was invalid because it incorrectly combined the total count of `POST /transactions` (which produce no events) and `POST /simulate` against the Kafka output. 
- By splitting the load tests into Tests A, B, and C, we can now accurately correlate the `POST /simulate` input with the Kafka event output and reconciliation results.
- **Intentional Failures:** The /simulate endpoint intentionally creates delayed events, missing ACCOUNTING_BOOKED events, and duplicate PAYMENT_RECEIVED events. `MISSING` and `DUPLICATE` are expected reconciliation results that prove the pipeline accurately reflects reality. They are not infrastructure failures.

## 5. DLQ Results
- The DLQ processes events that hard-fail to publish to Kafka (e.g., broker disconnects). Since expected simulation drops occur intentionally, they do not populate the DLQ. DLQ metrics should only be evaluated if infrastructure/connectivity issues occur.

## 6. Limitations
- The `POST /simulate` endpoint handles first-event generation synchronously. Under extreme load, this can still cause HTTP timeouts depending on database and Kafka availability. Batch inserts and asynchronous publishing may alleviate this bottleneck in the future.
