ALTER TABLE transaction_events
ALTER COLUMN received_at
TYPE TIMESTAMP WITH TIME ZONE
USING received_at AT TIME ZONE 'Asia/Kolkata';