CREATE TABLE rate_limits (
    key VARCHAR(500) PRIMARY KEY,
    count INTEGER NOT NULL DEFAULT 0,
    reset_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_rate_limits_reset_at ON rate_limits(reset_at);
