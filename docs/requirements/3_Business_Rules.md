BR-001

Every transaction must have a unique Transaction ID.

---

BR-002

Every event must belong to exactly one transaction.

---

BR-003

Each event must have a unique Event ID.

---

BR-004

A duplicate event must not change the final transaction state.

---

BR-005

A transaction is considered reconciled only when all required events have been received.

---

BR-006

Events may arrive in any order.

---

BR-007

Incomplete transactions should remain pending until either:

- Missing events arrive, or
- The reconciliation timeout expires.

---

BR-008

Events that repeatedly fail processing should be moved to the Dead Letter Queue.