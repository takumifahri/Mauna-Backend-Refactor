-- Create difficulty level enum
CREATE TYPE difficulty_level AS ENUM ('easy', 'medium', 'hard');

CREATE TABLE badges (
    id BIGSERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL UNIQUE,
    deskripsi TEXT,
    icon VARCHAR(255),
    level difficulty_level NOT NULL DEFAULT 'easy',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_badges_nama ON badges(nama);
CREATE INDEX idx_badges_level ON badges(level);

-- Add trigger for badges
CREATE TRIGGER trg_badges_set_updated_at
BEFORE UPDATE ON badges
FOR EACH ROW
EXECUTE FUNCTION set_timestamp_updated_at();