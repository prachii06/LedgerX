## Table: transactions

| Column | Type | Description |
|---------|------|-------------|
| transaction_id | UUID | Primary Key |
| status | VARCHAR | Current transaction state |
| created_at | TIMESTAMP | Creation time |
| updated_at | TIMESTAMP | Last update |

---

## Table: events

| Column | Type | Description |
|---------|------|-------------|
| event_id | UUID | Primary Key |
| transaction_id | UUID | Foreign Key |
| event_type | VARCHAR | Event type |
| source_service | VARCHAR | Event source |
| event_time | TIMESTAMP | Event timestamp |
| processed | BOOLEAN | Processing status |

---

## Table: reconciliation_records

| Column | Type | Description |
|---------|------|-------------|
| transaction_id | UUID | Primary & Foreign Key |
| status | VARCHAR | Reconciliation status |
| missing_events | JSONB | Expected but missing events |
| last_updated | TIMESTAMP | Last reconciliation |

---

## Table: failed_events

| Column | Type | Description |
|---------|------|-------------|
| id | UUID | Primary Key |
| event_id | UUID | Failed event |
| reason | TEXT | Failure reason |
| retry_count | INTEGER | Retry attempts |
| created_at | TIMESTAMP | Failure time |