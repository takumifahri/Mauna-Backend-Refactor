CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$
BEGIN
    CREATE TYPE user_role AS ENUM ('admin', 'user', 'moderator');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    CREATE TYPE user_tier AS ENUM ('bronze', 'silver', 'gold', 'diamond', 'platinum');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    CREATE TYPE kamus_category AS ENUM ('ALPHABET', 'NUMBERS', 'IMBUHAN', 'KOSAKATA');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    nama VARCHAR(255),
    telpon VARCHAR(255),
    role user_role NOT NULL DEFAULT 'user',
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
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

CREATE TABLE pending_registrations (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nama VARCHAR(255) NOT NULL,
    otp_hash VARCHAR(64) NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 0,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_pending_registrations_email ON pending_registrations(email);
CREATE INDEX ix_pending_registrations_expires_at ON pending_registrations(expires_at);

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
