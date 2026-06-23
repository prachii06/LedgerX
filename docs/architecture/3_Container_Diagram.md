# C4 Container Diagram

## Purpose

Shows the major software containers that make up LedgerX.

```text
+---------------------+
|   Dashboard (React) |
+----------+----------+
           |
           |
           v
+----------------------+
|     API Service      |
|      (Go + Gin)      |
+----------+-----------+
           |
     Publish Event
           |
           v
+----------------------+
|        Kafka         |
+----------+-----------+
           |
      Consume Event
           |
           v
+----------------------+
| Reconciliation Worker|
|        (Go)          |
+-----+-----------+----+
      |           |
      |           |
      v           v
 Redis Cache   PostgreSQL
```

## Containers

### Dashboard

Provides monitoring and search capabilities.

### API Service

Receives incoming transaction events.

### Kafka

Acts as the event broker.

### Worker

Processes and reconciles events.

### Redis

Stores temporary event state.

### PostgreSQL

Stores final transaction records.