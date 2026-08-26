CREATE TABLE reconciliation_results (
    id UUID PRIMARY KEY,
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    status VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    reconciled_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);