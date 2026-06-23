Indexes improve query performance.

## Transactions

- Primary Key: transaction_id
- Index on status
- Index on created_at

---

## Events

- Primary Key: event_id
- Index on transaction_id
- Index on event_type
- Index on source_service
- Index on event_time

---

## Reconciliation Records

- Primary Key: transaction_id
- Index on status

---

## Failed Events

- Primary Key: id
- Index on event_id
- Index on created_at