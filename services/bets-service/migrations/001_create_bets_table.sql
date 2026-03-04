CREATE TABLE IF NOT EXISTS price_bets (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    direction VARCHAR(10) NOT NULL,
    stake_amount DECIMAL(20, 8) NOT NULL,
    opened_price DECIMAL(20, 8) NOT NULL,
    resolved_price DECIMAL(20, 8) DEFAULT 0,
    status VARCHAR(20) NOT NULL,
    opened_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    resolved_at TIMESTAMP WITH TIME ZONE
);
