package entities

import "time"

type ShopItemType string

const (
    TypeStreakFreeze ShopItemType = "streak_freeze"
    TypeBadge        ShopItemType = "badge"
    TypeBoost        ShopItemType = "boost"
)

type ShopItem struct {
    ID          int64         `db:"id"`
    Name        string        `db:"name"`
    Description *string       `db:"description"`
    ItemType    ShopItemType  `db:"item_type"`
    XpCost      int           `db:"xp_cost"`
    Icon        *string       `db:"icon"`
    CreatedAt   time.Time     `db:"created_at"`
    UpdatedAt   *time.Time    `db:"updated_at"`
    DeletedAt   *time.Time    `db:"deleted_at"`
}

type Inventory struct {
    ID         int64      `db:"id"`
    UserID     int64      `db:"user_id"`
    ShopItemID int64      `db:"shop_item_id"`
    Quantity   int        `db:"quantity"`
    AcquiredAt time.Time  `db:"acquired_at"`
    CreatedAt  time.Time  `db:"created_at"`
    UpdatedAt  *time.Time `db:"updated_at"`
    DeletedAt  *time.Time `db:"deleted_at"`
}