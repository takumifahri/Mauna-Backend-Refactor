CREATE TABLE token_blacklist (
    id BIGSERIAL PRIMARY KEY,
    token VARCHAR(500) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    revoked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    reason VARCHAR(255)
);

CREATE INDEX ix_token_user_id ON token_blacklist(user_id);
CREATE INDEX ix_token_expires_at ON token_blacklist(expires_at);
CREATE INDEX ix_token_blacklist_token ON token_blacklist(token);