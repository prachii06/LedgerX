## ES-01: Duplicate Events

Scenario:
Same event is sent multiple times.

Handling:
- Detect using Event ID
- Ignore duplicates
- Do not update state

---

## ES-02: Missing Events

Scenario:
Payment event never arrives.

Handling:
- Mark transaction as PENDING
- After timeout → move to FAILED or REVIEW

---

## ES-03: Out-of-Order Events

Scenario:
Payment arrives before Order event.

Handling:
- Store in Redis
- Wait for missing dependency
- Reconcile later

---

## ES-04: Kafka Consumer Failure

Scenario:
Worker crashes while processing.

Handling:
- Kafka re-delivers event
- Processing remains idempotent

---

## ES-05: Database Failure

Scenario:
PostgreSQL temporarily unavailable.

Handling:
- Retry mechanism
- Events remain in Kafka until processed

---

## ES-06: Redis Failure

Scenario:
Redis is unavailable.

Handling:
- Fall back to PostgreSQL (slower path)
- System continues but with reduced performance

---

## ES-07: Worker Overload

Scenario:
Too many events arrive at once.

Handling:
- Kafka buffers load
- Horizontal scaling of workers (future improvement)