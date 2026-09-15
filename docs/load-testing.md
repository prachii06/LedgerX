# LedgerX Load Testing Report

## 1. Environment and Configuration
- **Test Environment:** Local machine (Windows).
- **k6 Version:** v2.2.0 (windows/amd64).
- **Database Configuration:** PostgreSQL running in Docker (`postgres` container, port 5433).
- **Kafka Configuration:** Kafka broker running in Docker.
- **Redis Configuration:** Redis running in Docker.

## 2. Test Scenarios

### Scenario A: Transactions Load Test (`transactions.js`)
- **Target:** `POST /transactions`
- **Goal:** Simulates typical external API client traffic.
- **Stages:** 10 VUs for 30s -> 50 VUs for 1m -> 100 VUs for 1m -> 0 VUs for 30s.
- **Duration:** 3 minutes
- **Requests generated:** 7868
- **Successful requests:** 7868
- **Failed requests:** 0
- **Average latency:** 9.45ms
- **p95 latency:** 18.55ms
- **p99 latency:** ~30ms (estimated from p95 and max)
- **Throughput:** 43.5 requests/second

### Scenario B: Simulation Load Test (`simulation.js`)
- **Target:** `POST /simulate`
- **Goal:** Stresses the internal bulk generation pipeline (internal transactions and events).
- **Stages:** 5 VUs for 30s -> 10 VUs for 1m -> 0 VUs for 30s.
- **Duration:** 2 minutes
- **Requests generated:** 111
- **Successful requests:** 93
- **Failed requests:** 18
- **Failure rate:** 16.21%
- **Average latency:** 4.12s
- **p95 latency:** 5.01s
- **Throughput:** ~0.9 requests/second

### Scenario C: Reconciliation Load Test (`reconciliation.js`)
- **Target:** `POST /transactions` followed by `GET /reconcile/:id`
- **Goal:** Simulates full processing loop and verifies eventual reconciliation status.
- **Stages:** 10 VUs for 30s -> 30 VUs for 1m -> 0 VUs for 30s.
- **Duration:** 2 minutes
- **Transactions generated:** 604
- **Reconciliation checks (successful HTTP):** 604
- **Failed HTTP requests:** 0
- **Average latency (overall):** 9.64ms
- **p95 latency:** 19.8ms

## 3. Database Integrity and Processing Metrics

After executing the load tests, the state of the PostgreSQL database was verified.

- **Transactions Created (Total):** 8989
- **Events Persisted in PostgreSQL:** 1512
- **Reconciliation Results Checked:** 9299
  - **MATCHED:** 0
  - **DUPLICATE:** 8
  - **MISSING:** 9291

## 4. Discovered Bottlenecks and Concurrency Issues

1. **Simulation API Timeout/Bottleneck:**
   - The `/simulate` endpoint severely degraded under moderate load (10 concurrent users). 
   - p95 latency reached 5 seconds, resulting in a 16.21% failure rate for HTTP requests crossing typical timeout thresholds.

2. **Severe Data Integrity Issue (Kafka / Event Processing):**
   - While 8989 transactions were created successfully, only 1512 transaction events made it to the PostgreSQL database.
   - The reconciliation engine checked over 9000 transactions, but marked 9291 as `MISSING` because the events never reached the database.
   - **Conclusion:** The event generation, Kafka producer, or Kafka consumer is silently dropping events under load, stalling, or crashing when bombarded with concurrent transaction creations. 

3. **Concurrency Violations (Duplicates):**
   - We observed 8 `DUPLICATE` reconciliation results, meaning the reconciliation worker processed the same transaction multiple times concurrently or received duplicate events from Kafka, bypassing idempotent checks.

4. **DLQ Architecture:**
   - The prompt suggested querying the DLQ table (`dlq_events`), however, it does not exist in the schema. This implies DLQ is either missing at the database level or stored elsewhere (like in a Kafka topic or Redis).

## 5. Grafana Monitoring

To monitor this in the future, the following PromQL queries/panels should be watched on Grafana:
- **HTTP:** `rate(http_requests_total[1m])`, `histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[1m]))`
- **Kafka:** `rate(kafka_events_produced_total[1m])`, `rate(kafka_events_consumed_total[1m])`, `rate(producer_errors_total[1m])`
- **Reconciliation:** `sum by (status) (reconciliations_total)`
- **DLQ:** `rate(event_dlq_total[1m])`
