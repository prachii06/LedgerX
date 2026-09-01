ALTER TABLE reconciliation_results
ALTER COLUMN reconciled_at
TYPE TIMESTAMPTZ
USING reconciled_at AT TIME ZONE 'UTC';