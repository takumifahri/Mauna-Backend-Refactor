CREATE TABLE sublevel (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    tujuan TEXT,
    level_id BIGINT NOT NULL REFERENCES level(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_sublevel_level_id ON sublevel(level_id);
CREATE INDEX idx_sublevel_name ON sublevel(name);

CREATE TRIGGER trg_sublevel_set_updated_at
BEFORE UPDATE ON sublevel
FOR EACH ROW
EXECUTE FUNCTION set_timestamp_updated_at();