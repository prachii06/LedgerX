ALTER TABLE reconciliation_results
ALTER COLUMN reconciled_at
TYPE TIMESTAMP
USING reconciled_at AT TIME ZONE 'UTC';