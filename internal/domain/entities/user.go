package entities

import "time"

type UserRole string
type UserTier string

const (
    RoleAdmin     UserRole = "admin"
    RoleUser      UserRole = "user"
    RoleModerator UserRole = "moderator"
)

const (
    TierBronze    UserTier = "bronze"
    TierSilver    UserTier = "silver"
    TierGold      UserTier = "gold"
    TierDiamond   UserTier = "diamond"
    TierPlatinum  UserTier = "platinum"
)

type User struct {
    ID                    int64      `db:"id"`
    UniqueID              string     `db:"unique_id"`
    Username              string     `db:"username"`
    Email                 string     `db:"email"`
    Password              string     `db:"password"`
    Nama                  *string    `db:"nama"`
    Telpon                *string    `db:"telpon"`
    Role                  UserRole   `db:"role"`
    IsActive              bool       `db:"is_active"`
    IsVerified            bool       `db:"is_verified"`
    Avatar                *string    `db:"avatar"`
    Bio                   *string    `db:"bio"`
    TotalBadges           int        `db:"total_badges"`
    AvatarURL             *string    `db:"avatar_url"`
    CurrentStreak         int        `db:"current_streak"`
    LongestStreak         int        `db:"longest_streak"`
    LastActivityDate      *time.Time `db:"last_activity_date"`
    StreakFreezeCount     int        `db:"streak_freeze_count"`
    WeeklyXp              int        `db:"weekly_xp"`
    Tier                  UserTier   `db:"tier"`
    TotalXp               int        `db:"total_xp"`
    TotalQuizzesCompleted int        `db:"total_quizzes_completed"`
    TotalPoints           int        `db:"total_points"`
    CreatedAt             time.Time  `db:"created_at"`
    UpdatedAt             *time.Time `db:"updated_at"`
    DeletedAt             *time.Time `db:"deleted_at"`
    LastLogin             *time.Time `db:"last_login"`
}