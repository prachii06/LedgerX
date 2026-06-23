The following assumptions apply to Version 1 of LedgerX.

- Each transaction has a globally unique Transaction ID.
- Each event has a globally unique Event ID.
- Events are immutable after creation.
- Events are processed at least once.
- Duplicate events are expected.
- Events may arrive out of order.
- Some events may never arrive.
- PostgreSQL is the source of truth.
- Redis stores temporary reconciliation data.
- Kafka guarantees ordering within a partition.


