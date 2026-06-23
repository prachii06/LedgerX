# Event Flow

## Normal Flow

Transaction Created

↓

API receives event

↓

Kafka stores event

↓

Worker consumes event

↓

Duplicate check

↓

Reconciliation

↓

Save transaction

↓

Dashboard displays status

---

## Failure Flow

Event

↓

Worker failure

↓

Kafka retries

↓

Worker restarts

↓

Process resumes

↓

Transaction updated