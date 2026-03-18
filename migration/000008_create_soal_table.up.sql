-- Create soal type enum
CREATE TYPE soal_type AS ENUM ('TEBAK_GAMBAR', 'OPEN_CAMERA', 'PILIHAN_GANDA', 'MATEMATIKA');

CREATE TABLE soal (
    id BIGSERIAL PRIMARY KEY,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    video_url TEXT,
    image_url TEXT,
    dictionary_id BIGINT NOT NULL REFERENCES kamus(id) ON DELETE CASCADE,
    point_gamifikasi INTEGER NOT NULL DEFAULT 10,
    sublevel_id BIGINT NOT NULL REFERENCES sublevel(id) ON DELETE CASCADE,
    categories soal_type NOT NULL DEFAULT 'OPEN_CAMERA',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_soal_dictionary_id ON soal(dictionary_id);
CREATE INDEX idx_soal_sublevel_id ON soal(sublevel_id);
CREATE INDEX idx_soal_categories ON soal(categories);

CREATE TRIGGER trg_soal_set_updated_at
BEFORE UPDATE ON soal
FOR EACH ROW
EXECUTE FUNCTION set_timestamp_updated_at();