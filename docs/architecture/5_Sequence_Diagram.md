# Sequence Diagrams

## Event Processing

```text
Client
  |
  | POST /events
  |
  v
API Service
  |
  | Publish
  |
  v
Kafka
  |
  | Consume
  |
  v
Worker
  |
  | Check Duplicate
  |
  v
Redis
  |
  | Reconcile
  |
  v
PostgreSQL
```

## Result

The transaction state is updated after reconciliation.