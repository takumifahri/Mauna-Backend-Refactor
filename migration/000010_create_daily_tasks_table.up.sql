CREATE TABLE daily_task (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE,
    completed_sublevels INTEGER NOT NULL DEFAULT 0,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    last_update TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_daily_task_user_id ON daily_task(user_id);
CREATE INDEX idx_daily_task_date ON daily_task(date);

CREATE TRIGGER trg_daily_task_set_updated_at
BEFORE UPDATE ON daily_task
FOR EACH ROW
EXECUTE FUNCTION set_timestamp_updated_at();