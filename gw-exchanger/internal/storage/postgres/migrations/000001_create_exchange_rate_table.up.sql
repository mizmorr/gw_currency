CREATE TABLE exchange_rate (
    currency_code VARCHAR(3) PRIMARY KEY,
    rate DOUBLE PRECISION NOT NULL
);

INSERT INTO exchange_rate(currency_code, rate) VALUES
('USD',101.9146),
('RUB',1.0),
('EUR',105.0464)

