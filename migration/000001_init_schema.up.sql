CREATE TYPE user_role AS ENUM ('admin', 'user', 'moderator');
CREATE TYPE user_tier AS ENUM ('bronze', 'silver', 'gold', 'diamond', 'platinum');
CREATE TYPE kamus_category AS ENUM ('ALPHABET', 'NUMBERS', 'IMBUHAN', 'KOSAKATA');

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    unique_id VARCHAR(50) UNIQUE,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    nama VARCHAR(255),
    telpon VARCHAR(255),
    role user_role NOT NULL DEFAULT 'user',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    avatar VARCHAR(255),
    bio TEXT,
    total_badges INTEGER NOT NULL DEFAULT 0,
    avatar_url VARCHAR(255),
    current_streak INTEGER NOT NULL DEFAULT 0,
    longest_streak INTEGER NOT NULL DEFAULT 0,
    last_activity_date DATE,
    streak_freeze_count INTEGER NOT NULL DEFAULT 0,
    weekly_xp INTEGER NOT NULL DEFAULT 0,
    tier user_tier NOT NULL DEFAULT 'bronze',
    total_xp INTEGER NOT NULL DEFAULT 0,
    total_quizzes_completed INTEGER NOT NULL DEFAULT 0,
    total_points INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    last_login TIMESTAMPTZ
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_unique_id ON users(unique_id);

CREATE TABLE kamus (
    id BIGSERIAL PRIMARY KEY,
    word_text VARCHAR(255) NOT NULL UNIQUE,
    definition TEXT NOT NULL,
    category kamus_category NOT NULL DEFAULT 'ALPHABET',
    image_url_ref VARCHAR(255),
    video_url VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_kamus_word_text ON kamus(word_text);
CREATE INDEX idx_kamus_category ON kamus(category);

CREATE OR REPLACE FUNCTION set_timestamp_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_timestamp_updated_at();

CREATE TRIGGER trg_kamus_set_updated_at
BEFORE UPDATE ON kamus
FOR EACH ROW
EXECUTE FUNCTION set_timestamp_updated_at();

CREATE OR REPLACE FUNCTION set_users_unique_id_after_insert()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.unique_id IS NULL OR NEW.unique_id = '' THEN
    UPDATE users
    SET unique_id = 'USR-' || LPAD(NEW.id::TEXT, 5, '0')
    WHERE id = NEW.id;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_set_unique_id_after_insert
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION set_users_unique_id_after_insert();