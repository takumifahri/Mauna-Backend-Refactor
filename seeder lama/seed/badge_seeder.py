from typing import List, Dict, Any
from sqlalchemy.orm import Session
from datetime import datetime

from ...models.badges import Badge, DificultyLevel
from ...config.hash import hash_password
from ..seeder import BaseSeeder

class BadgeSeeder(BaseSeeder):
    """Seed initial badges into the database"""
    
    def run(self):
        """Run badge seeding"""
        try:
            print("🌱 Seeding badges...")
            
            badges_data = [
                {
                    "nama": "First Steps",
                    "deskripsi": "Complete your first lesson",
                    "icon": "badges/First_Steps.png",
                    "level": DificultyLevel.EASY
                },
                {
                    "nama": "Alphabet Master",
                    "deskripsi": "Master all alphabet signs",
                    "icon": "badges/Alphabet_Master.png",
                    "level": DificultyLevel.MEDIUM
                },
                {
                    "nama": "Number Expert",
                    "deskripsi": "Perfect number signs 0-9",
                    "icon": "badges/Numbers_Expert.png",
                    "level": DificultyLevel.MEDIUM
                },
                {
                    "nama": "Conversation Pro",
                    "deskripsi": "Complete advanced conversations",
                    "icon": "badges/Conversation.png",
                    "level": DificultyLevel.HARD
                }
            ]
            
            created_count = 0
            for badge_data in badges_data:
                existing_badge = self.db.query(Badge).filter(
                    Badge.nama == badge_data["nama"]
                ).first()
                
                if not existing_badge:
                    badge = Badge(**badge_data)
                    self.db.add(badge)
                    created_count += 1
                    print(f"  ✅ Created badge: {badge_data['nama']}")
                else:
                    print(f"  ⚠️ Badge already exists: {badge_data['nama']}")
            
            self.db.commit()
            print(f"✅ Badge seeding completed. Created {created_count} new badges.")
            
        except Exception as e:
            self.db.rollback()
            raise Exception(f"Badge seeding failed: {e}")