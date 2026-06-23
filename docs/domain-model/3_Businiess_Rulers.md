## BR-01

Every transaction must have a unique Transaction ID.

---

## BR-02

Every event must have a unique Event ID.

---

## BR-03

An event belongs to exactly one transaction.

---

## BR-04

A transaction may contain multiple events.

---

## BR-05

Duplicate events must never change the final transaction state.

---

## BR-06

Transactions become RECONCILED only when all required events are received.

---

## BR-07

Events may arrive in any order.

The system must still reconcile correctly.

---

## BR-08

Missing events should trigger a timeout.

---

## BR-09

Failed events should be moved to the Dead Letter Queue after retry attempts are exhausted.

---

## BR-10

The final transaction state must be stored permanently in PostgreSQL.