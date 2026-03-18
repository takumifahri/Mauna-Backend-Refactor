CREATE TABLE level (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    description TEXT,
    tujuan TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_level_name ON level(name);

CREATE TRIGGER trg_level_set_updated_at
BEFORE UPDATE ON level
FOR EACH ROW
EXECUTE FUNCTION set_timestamp_updated_at();