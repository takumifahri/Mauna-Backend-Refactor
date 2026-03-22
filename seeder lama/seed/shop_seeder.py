from typing import List, Dict, Any
from sqlalchemy.orm import Session

try:
    from src.models.shop_items import ShopItem, ShopItemType
    from src.database.seeder import BaseSeeder
except ImportError:
    from ...models.shop_items import ShopItem, ShopItemType
    from ..seeder import BaseSeeder

class ShopSeeder(BaseSeeder):
    """Seeder untuk Shop Items"""

    def __init__(self):
        super().__init__()

    def run(self):
        try:
            print("🌱 Starting shop item seeding...")
            self.seed_shop_items()
            print("✅ Shop item seeding finished successfully!")
        except Exception as e:
            self.db.rollback()
            raise Exception(f"Shop item seeding failed: {e}")

    def seed_shop_items(self):
        print("🛒 Seeding Shop Items...")
        items_data = [
            {
                "name": "Streak Freeze",
                "description": "Melindungi streak harian agar tidak hilang.",
                "item_type": ShopItemType.STREAK_FREEZE,
                "xp_cost": 100,
                "icon": "streak_freeze.png"
            },
            {
                "name": "Badge Emas",
                "description": "Badge spesial untuk pencapaian emas.",
                "item_type": ShopItemType.BADGE,
                "xp_cost": 250,
                "icon": "badge_emas.png"
            },
            {
                "name": "XP Boost",
                "description": "Dapatkan XP tambahan selama 1 jam.",
                "item_type": ShopItemType.BOOST,
                "xp_cost": 150,
                "icon": "xp_boost.png"
            },
        ]

        created_count = 0
        for data in items_data:
            existing = self.db.query(ShopItem).filter(ShopItem.name == data["name"]).first()
            if not existing:
                item = ShopItem(**data)
                self.db.add(item)
                created_count += 1
                print(f"  ✅ Created shop item: {data['name']}")
            else:
                print(f"  ⚠️ Shop item already exists: {data['name']}")

        self.db.commit()
        print(f"✅ Shop item seeding completed. Created {created_count} items.")

if __name__ == "__main__":
    ShopSeeder().run()