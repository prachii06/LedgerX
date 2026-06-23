## Flow 1: Normal Successful Flow

Transaction Created → Order Event → Payment Event → Accounting Event → Reconciled

System Path:

Client → API → Kafka → Worker → Redis → PostgreSQL

Final State:
RECONCILED

---

## Flow 2: Out-of-Order Flow

Payment Event arrives first

1. Payment Event → stored in Redis
2. Order Event arrives later
3. Reconciliation triggered
4. Transaction marked RECONCILED

---

## Flow 3: Duplicate Events

Event A received twice

1. First event processed
2. Second event detected as duplicate
3. Ignored

---

## Flow 4: Missing Event

Order + Payment arrive but Accounting missing

1. System waits
2. Timeout triggers
3. Transaction marked PENDING_REVIEW

---

## Flow 5: Failure Recovery

Worker crashes mid-processing

1. Kafka re-delivers event
2. Worker resumes processing
3. No data loss due to idempotency