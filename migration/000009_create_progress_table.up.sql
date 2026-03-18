-- Create progress status enum
CREATE TYPE progress_status AS ENUM ('not_started', 'in_progress', 'completed', 'failed');

CREATE TABLE progress (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sublevel_id BIGINT NOT NULL REFERENCES sublevel(id) ON DELETE CASCADE,
    status progress_status NOT NULL DEFAULT 'not_started',
    total_questions INTEGER NOT NULL DEFAULT 0,
    correct_answers INTEGER NOT NULL DEFAULT 0,
    score INTEGER NOT NULL DEFAULT 0,
    stars INTEGER NOT NULL DEFAULT 0,
    completion_percentage INTEGER NOT NULL DEFAULT 0,
    attempts INTEGER NOT NULL DEFAULT 0,
    best_score INTEGER NOT NULL DEFAULT 0,
    best_stars INTEGER NOT NULL DEFAULT 0,
    is_unlocked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    last_attempt TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    CONSTRAINT uq_user_sublevel_progress UNIQUE(user_id, sublevel_id)
);

CREATE INDEX idx_progress_user_id ON progress(user_id);
CREATE INDEX idx_progress_sublevel_id ON progress(sublevel_id);
CREATE INDEX idx_progress_status ON progress(status);

CREATE TRIGGER trg_progress_set_updated_at
BEFORE UPDATE ON progress
FOR EACH ROW
EXECUTE FUNCTION set_timestamp_updated_at();