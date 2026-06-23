# Glossary

## Transaction

A business operation that represents a financial activity, such as an online payment.

---

## Event

A message describing something that has happened within the system.

Examples:

- Order Created
- Payment Received
- Accounting Recorded

---

## Reconciliation

The process of matching related events from different services to verify that a transaction has been processed correctly.

---

## Idempotency

The ability to process the same event multiple times without changing the final result.

---

## Dead Letter Queue (DLQ)

A queue where events that cannot be processed successfully are stored for later investigation or replay.

---

## Worker

A background service that consumes events from Kafka and performs reconciliation.

---

## Event Stream

A continuous flow of events published by different services.

---

## Operations Dashboard

A web interface used to monitor transaction processing, worker status, and system health.