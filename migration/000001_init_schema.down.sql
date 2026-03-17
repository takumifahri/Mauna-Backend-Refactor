DROP TRIGGER IF EXISTS trg_users_set_unique_id_after_insert ON users;
DROP FUNCTION IF EXISTS set_users_unique_id_after_insert();

DROP TRIGGER IF EXISTS trg_users_set_updated_at ON users;
DROP TRIGGER IF EXISTS trg_kamus_set_updated_at ON kamus;
DROP FUNCTION IF EXISTS set_timestamp_updated_at();

DROP TABLE IF EXISTS kamus;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS kamus_category;
DROP TYPE IF EXISTS user_tier;
DROP TYPE IF EXISTS user_role;