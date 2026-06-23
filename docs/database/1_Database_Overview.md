## Purpose

LedgerX uses PostgreSQL as its primary database.

It stores the final state of transactions, processed events, reconciliation results, and failed events. PostgreSQL acts as the source of truth for all persistent data.

Redis is used only for temporary processing state and caching.

---

## Database Responsibilities

The database is responsible for:

- Storing transaction records
- Storing processed events
- Maintaining reconciliation status
- Recording failed events
- Supporting dashboard queries
- Preserving historical data

---

## Database Choice

**PostgreSQL**

Reasons:

- ACID-compliant
- Reliable persistence
- Excellent indexing support
- Strong SQL capabilities
- Mature ecosystem