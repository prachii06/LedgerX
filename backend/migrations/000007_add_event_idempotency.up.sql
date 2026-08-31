ALTER TABLE transaction_events
ADD CONSTRAINT unique_event_id UNIQUE (id);