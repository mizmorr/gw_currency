-- migrations/003_wallet_balances_table.up.sql

CREATE TABLE wallet_balances (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL,
    currency VARCHAR(10) NOT NULL,
    balance DECIMAL(20, 2) DEFAULT 0.00,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(wallet_id, currency),
    FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION update_balance_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_balance_updated_at
BEFORE UPDATE ON wallet_balances
FOR EACH ROW
EXECUTE FUNCTION update_balance_timestamp();
