## 1. Transaction

Represents a business transaction.

Examples:
- Online payment
- Order payment
- Refund

Attributes:

- Transaction ID
- Current Status
- Created Time
- Updated Time
- Expected Events
- Received Events

---

## 2. Event

Represents something that happened within another system.

Examples:

- Order Created
- Payment Completed
- Accounting Recorded

Attributes:

- Event ID
- Transaction ID
- Event Type
- Timestamp
- Source Service

---

## 3. Reconciliation Record

Represents the reconciliation result for a transaction.

Attributes:

- Transaction ID
- Reconciliation Status
- Missing Events
- Last Updated

---

## 4. Worker

Background processor responsible for consuming events.

Responsibilities:

- Consume Kafka events
- Detect duplicates
- Update reconciliation state
- Store final transaction

---

## 5. Dead Letter Event

Represents an event that could not be processed.

Attributes:

- Event ID
- Failure Reason
- Retry Count
- Timestamp