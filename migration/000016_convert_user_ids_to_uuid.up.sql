CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'users'
          AND column_name = 'id'
          AND data_type = 'bigint'
    ) THEN
        ALTER TABLE users ADD COLUMN new_id UUID NOT NULL DEFAULT gen_random_uuid();

        IF EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_name = 'users'
              AND column_name = 'unique_id'
        ) THEN
            UPDATE users
            SET new_id = unique_id::uuid
            WHERE unique_id IS NOT NULL
              AND unique_id ~* '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$';
        END IF;

        ALTER TABLE user_badges DROP CONSTRAINT IF EXISTS user_badges_user_id_fkey;
        ALTER TABLE progress DROP CONSTRAINT IF EXISTS progress_user_id_fkey;
        ALTER TABLE daily_task DROP CONSTRAINT IF EXISTS daily_task_user_id_fkey;
        ALTER TABLE inventory DROP CONSTRAINT IF EXISTS inventory_user_id_fkey;
        ALTER TABLE password_reset_tokens DROP CONSTRAINT IF EXISTS password_reset_tokens_user_id_fkey;

        ALTER TABLE user_badges DROP CONSTRAINT IF EXISTS user_badges_user_id_badge_id_key;
        ALTER TABLE progress DROP CONSTRAINT IF EXISTS uq_user_sublevel_progress;

        ALTER TABLE user_badges ADD COLUMN new_user_id UUID;
        UPDATE user_badges ub SET new_user_id = u.new_id FROM users u WHERE ub.user_id = u.id;
        ALTER TABLE user_badges ALTER COLUMN new_user_id SET NOT NULL;
        ALTER TABLE user_badges DROP COLUMN user_id;
        ALTER TABLE user_badges RENAME COLUMN new_user_id TO user_id;

        ALTER TABLE progress ADD COLUMN new_user_id UUID;
        UPDATE progress p SET new_user_id = u.new_id FROM users u WHERE p.user_id = u.id;
        ALTER TABLE progress ALTER COLUMN new_user_id SET NOT NULL;
        ALTER TABLE progress DROP COLUMN user_id;
        ALTER TABLE progress RENAME COLUMN new_user_id TO user_id;

        ALTER TABLE daily_task ADD COLUMN new_user_id UUID;
        UPDATE daily_task dt SET new_user_id = u.new_id FROM users u WHERE dt.user_id = u.id;
        ALTER TABLE daily_task ALTER COLUMN new_user_id SET NOT NULL;
        ALTER TABLE daily_task DROP COLUMN user_id;
        ALTER TABLE daily_task RENAME COLUMN new_user_id TO user_id;

        ALTER TABLE inventory ADD COLUMN new_user_id UUID;
        UPDATE inventory i SET new_user_id = u.new_id FROM users u WHERE i.user_id = u.id;
        ALTER TABLE inventory ALTER COLUMN new_user_id SET NOT NULL;
        ALTER TABLE inventory DROP COLUMN user_id;
        ALTER TABLE inventory RENAME COLUMN new_user_id TO user_id;

        ALTER TABLE password_reset_tokens ADD COLUMN new_user_id UUID;
        UPDATE password_reset_tokens prt SET new_user_id = u.new_id FROM users u WHERE prt.user_id = u.id;
        ALTER TABLE password_reset_tokens ALTER COLUMN new_user_id SET NOT NULL;
        ALTER TABLE password_reset_tokens DROP COLUMN user_id;
        ALTER TABLE password_reset_tokens RENAME COLUMN new_user_id TO user_id;

        ALTER TABLE token_blacklist ADD COLUMN new_user_id UUID;
        UPDATE token_blacklist tb SET new_user_id = u.new_id FROM users u WHERE tb.user_id = u.id;
        ALTER TABLE token_blacklist ALTER COLUMN new_user_id SET NOT NULL;
        ALTER TABLE token_blacklist DROP COLUMN user_id;
        ALTER TABLE token_blacklist RENAME COLUMN new_user_id TO user_id;

        ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pkey;
        ALTER TABLE users RENAME COLUMN id TO legacy_id;
        ALTER TABLE users RENAME COLUMN new_id TO id;
        ALTER TABLE users ADD PRIMARY KEY (id);
        ALTER TABLE users DROP COLUMN legacy_id;
        ALTER TABLE users DROP COLUMN IF EXISTS unique_id;

        DROP INDEX IF EXISTS idx_users_unique_id;
        CREATE INDEX IF NOT EXISTS idx_user_badges_user_id ON user_badges(user_id);
        CREATE INDEX IF NOT EXISTS idx_user_badges_user_earned ON user_badges(user_id, earned_at);
        CREATE INDEX IF NOT EXISTS idx_progress_user_id ON progress(user_id);
        CREATE INDEX IF NOT EXISTS idx_daily_task_user_id ON daily_task(user_id);
        CREATE INDEX IF NOT EXISTS idx_inventory_user_id ON inventory(user_id);
        CREATE INDEX IF NOT EXISTS ix_token_user_id ON token_blacklist(user_id);
        CREATE INDEX IF NOT EXISTS ix_password_reset_tokens_user_id ON password_reset_tokens(user_id);

        ALTER TABLE user_badges ADD CONSTRAINT user_badges_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
        ALTER TABLE progress ADD CONSTRAINT progress_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
        ALTER TABLE daily_task ADD CONSTRAINT daily_task_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
        ALTER TABLE inventory ADD CONSTRAINT inventory_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
        ALTER TABLE password_reset_tokens ADD CONSTRAINT password_reset_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

        ALTER TABLE user_badges ADD CONSTRAINT user_badges_user_id_badge_id_key UNIQUE(user_id, badge_id);
        ALTER TABLE progress ADD CONSTRAINT uq_user_sublevel_progress UNIQUE(user_id, sublevel_id);
    ELSE
        ALTER TABLE users DROP COLUMN IF EXISTS unique_id;
        DROP INDEX IF EXISTS idx_users_unique_id;
    END IF;
END $$;
