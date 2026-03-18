CREATE TYPE shop_item_type AS ENUM ('streak_freeze', 'badge', 'boost');

CREATE TABLE shop_items (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    item_type shop_item_type NOT NULL,
    xp_cost INTEGER NOT NULL,
    icon VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_shop_items_item_type ON shop_items(item_type);

CREATE TRIGGER trg_shop_items_set_updated_at
BEFORE UPDATE ON shop_items
FOR EACH ROW
EXECUTE FUNCTION set_timestamp_updated_at();