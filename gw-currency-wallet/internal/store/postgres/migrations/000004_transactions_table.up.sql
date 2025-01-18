-- migrations/004_transactions_table.up.sql

CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL,
    currency VARCHAR(10) NOT NULL,
    amount DECIMAL(20, 2) NOT NULL,
    type VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE
);

