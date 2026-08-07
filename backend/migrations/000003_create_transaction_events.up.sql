CREATE TABLE transaction_events (
    id UUID PRIMARY KEY,

    transaction_id UUID NOT NULL,

    source VARCHAR(50) NOT NULL,

    event_type VARCHAR(100) NOT NULL,

    payload JSONB,

    received_at TIMESTAMP NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_transaction
        FOREIGN KEY (transaction_id)
        REFERENCES transactions(id)
        ON DELETE CASCADE
);