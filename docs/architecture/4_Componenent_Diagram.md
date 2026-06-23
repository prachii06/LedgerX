# Component Diagram

## API Service Components

```text
HTTP API
    │
    ▼
Validation
    │
    ▼
Event Publisher
    │
    ▼
Kafka
```

---

## Worker Components

```text
Kafka Consumer
       │
       ▼
Duplicate Checker
       │
       ▼
Reconciliation Engine
       │
       ▼
State Manager
       │
       ▼
Database Writer
```

## Responsibilities

### Validation

Checks event format.

### Event Publisher

Publishes validated events.

### Duplicate Checker

Ensures idempotent processing.

### Reconciliation Engine

Matches related events.

### State Manager

Determines transaction status.

### Database Writer

Persists final transaction state.