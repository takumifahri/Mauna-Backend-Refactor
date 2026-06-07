DO $$
BEGIN
    RAISE EXCEPTION 'Converting users.id from UUID back to BIGSERIAL is not safely reversible without losing ID mappings.';
END $$;
