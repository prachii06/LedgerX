CREATE TABLE transactions (
    id UUID PRIMARY KEY,

    external_id VARCHAR(100) UNIQUE NOT NULL,

    amount NUMERIC(12,2) NOT NULL,

    currency VARCHAR(3) NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',

    created_at TIMESTAMP DEFAULT NOW(),

    updated_at TIMESTAMP DEFAULT NOW()
);