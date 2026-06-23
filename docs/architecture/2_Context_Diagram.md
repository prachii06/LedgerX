## Purpose

The Context Diagram shows how LedgerX interacts with external actors and systems.

```text
                  +----------------------+
                  | Transaction Generator|
                  +----------+-----------+
                             |
                             | Sends Events
                             |
                             v
                     +---------------+
                     |    LedgerX    |
                     +---------------+
                             |
            +----------------+----------------+
            |                |                |
            v                v                v
      PostgreSQL         Redis           Operations
                                            Dashboard
```

## External Actors

### Transaction Generator

Produces transaction events.

Examples:

- Order Service
- Payment Service
- Accounting Service

### Operations Dashboard

Displays transaction status, metrics, and logs.