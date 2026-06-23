## Migration Tool

Database schema changes will be managed using version-controlled SQL migrations.

Each migration should:

- Be incremental
- Be reversible where possible
- Be committed to version control

---

## Naming Convention

Example:

001_create_transactions.sql

002_create_events.sql

003_create_indexes.sql

004_add_failed_events.sql

---

## Benefits

- Consistent database setup
- Easy rollback
- Team collaboration
- Automated deployment