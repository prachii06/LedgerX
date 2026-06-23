## UC-01: Create Transaction Event

Actor: Transaction Generator

Flow:
1. A transaction event is generated.
2. The event is sent to the LedgerX API.
3. API validates the event.
4. Event is published to Kafka.

Result:
- Event enters the processing pipeline.

---

## UC-02: Process Event

Actor: Reconciliation Worker

Flow:
1. Worker consumes event from Kafka.
2. Event is checked for duplicates.
3. Event is stored in Redis (if needed).
4. Matching logic is applied.
5. Transaction state is updated in PostgreSQL.

Result:
- Event is either reconciled or marked pending.

---

## UC-03: Reconcile Transaction

Actor: Reconciliation Worker

Flow:
1. Worker receives related events for a transaction.
2. System checks if all required events exist:
   - Order Event
   - Payment Event
   - Accounting Event
3. If all exist → mark transaction as RECONCILED.

Result:
- Transaction status updated.

---

## UC-04: Handle Out-of-Order Event

Actor: Reconciliation Worker

Flow:
1. Event arrives before its dependent event.
2. Event is stored in Redis as pending.
3. When missing event arrives → reconciliation is triggered.

Result:
- System eventually reconciles correctly.

---

## UC-05: Detect Duplicate Event

Actor: Reconciliation Worker

Flow:
1. Event ID is checked.
2. If already processed:
   - Ignore event
   - Do not update state

Result:
- Idempotent processing ensured.

---

## UC-06: Handle Missing Event Timeout

Actor: Reconciliation Worker (Background Job)

Flow:
1. System checks incomplete transactions.
2. If timeout exceeded:
   - Mark transaction as FAILED or PENDING_REVIEW

Result:
- System avoids infinite waiting.

---

## UC-07: View Transaction in Dashboard

Actor: User (Engineer / Operator)

Flow:
1. User opens dashboard.
2. Searches transaction ID.
3. System fetches data from PostgreSQL.
4. Timeline is displayed.

Result:
- User sees full transaction lifecycle.

---

## UC-08: View Live Event Stream

Actor: User

Flow:
1. Dashboard subscribes to event stream.
2. Events are pushed in real-time.
3. UI updates automatically.

Result:
- Live observability of system.