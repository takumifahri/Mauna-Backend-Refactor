package seeder

import (
    "database/sql"
    "fmt"
    "github.com/jmoiron/sqlx"
    "time"
)

type UserBadgeSeeder struct {
    db *sqlx.DB
}

func NewUserBadgeSeeder(db *sqlx.DB) *UserBadgeSeeder {
    return &UserBadgeSeeder{db: db}
}

func (s *UserBadgeSeeder) Name() string {
    return "UserBadgeSeeder"
}

type userBadgeMapping struct {
    username string
    badgeIDs []int64
}

func (s *UserBadgeSeeder) Run() error {
    PrintInfo("🏆 Seeding user badges...")

    tx, err := s.db.Beginx()
    if err != nil {
        PrintError(fmt.Sprintf("Failed to start transaction: %v", err))
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    var userCount, badgeCount int
    if err := tx.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount); err != nil || userCount == 0 {
        PrintWarning("No users found. Please run UserSeeder first.")
        return nil
    }

    if err := tx.QueryRow("SELECT COUNT(*) FROM badges").Scan(&badgeCount); err != nil || badgeCount == 0 {
        PrintWarning("No badges found. Please run BadgeSeeder first.")
        return nil
    }

    mappings := []userBadgeMapping{
        {username: "admin", badgeIDs: []int64{1, 2, 3, 4}},
        {username: "moderator", badgeIDs: []int64{1, 2, 3}},
        {username: "johndoe", badgeIDs: []int64{1}},
        {username: "janedoe", badgeIDs: []int64{1, 2}},
    }

    assignedCount := 0
    for _, mapping := range mappings {
        var userID int64
        err := tx.QueryRow("SELECT id FROM users WHERE username = $1", mapping.username).Scan(&userID)

        if err == sql.ErrNoRows {
            PrintWarning(fmt.Sprintf("User not found: %s", mapping.username))
            continue
        }
        if err != nil {
            PrintError(fmt.Sprintf("Database error: %v", err))
            return err
        }

        for _, badgeID := range mapping.badgeIDs {
            var existingID sql.NullInt64
            err := tx.QueryRow(
                "SELECT id FROM user_badges WHERE user_id = $1 AND badge_id = $2",
                userID, badgeID,
            ).Scan(&existingID)

            if err == nil && existingID.Valid {
                badgeName := s.getBadgeNameTx(tx, badgeID)
                PrintWarning(fmt.Sprintf("%s already has '%s'", mapping.username, badgeName))
                continue
            }

            if err != nil && err != sql.ErrNoRows {
                PrintError(fmt.Sprintf("Database error: %v", err))
                return err
            }

            badgeName := s.getBadgeNameTx(tx, badgeID)
            _, err = tx.Exec(
                "INSERT INTO user_badges (user_id, badge_id, earned_at) VALUES ($1, $2, $3)",
                userID, badgeID, time.Now(),
            )

            if err != nil {
                PrintError(fmt.Sprintf("Failed to award badge: %v", err))
                return err
            }

            PrintSuccess(fmt.Sprintf("%s earned '%s' (ID: %d)", mapping.username, badgeName, badgeID))
            assignedCount++
        }

        totalBadges := s.getUserBadgeCountTx(tx, userID)
        _, err = tx.Exec("UPDATE users SET total_badges = $1 WHERE id = $2", totalBadges, userID)
        if err != nil {
            PrintError(fmt.Sprintf("Failed to update user total_badges: %v", err))
            return err
        }
    }

    if err := tx.Commit(); err != nil {
        PrintError(fmt.Sprintf("Failed to commit transaction: %v", err))
        return err
    }

    PrintSuccess(fmt.Sprintf("Assigned %d badges to users.", assignedCount))
    PrintSuccess("User badges seeding completed.")
    return nil
}

// getBadgeNameTx retrieves badge name by ID from transaction
func (s *UserBadgeSeeder) getBadgeNameTx(tx *sqlx.Tx, badgeID int64) string {
    var nama string
    if err := tx.QueryRow("SELECT nama FROM badges WHERE id = $1", badgeID).Scan(&nama); err != nil {
        return fmt.Sprintf("Badge ID %d", badgeID)
    }
    return nama
}

// getUserBadgeCountTx returns total badge count for a user (transaction)
func (s *UserBadgeSeeder) getUserBadgeCountTx(tx *sqlx.Tx, userID int64) int {
    var count int
    tx.QueryRow("SELECT COUNT(*) FROM user_badges WHERE user_id = $1", userID).Scan(&count)
    return count
}

// getUserBadgeCountDB returns total badge count for a user (direct DB)
func (s *UserBadgeSeeder) getUserBadgeCountDB(userID int64) int {
    var count int
    s.db.QueryRow("SELECT COUNT(*) FROM user_badges WHERE user_id = $1", userID).Scan(&count)
    return count
}

// AssignBadgeToUser manually assign badge to user by ID
func (s *UserBadgeSeeder) AssignBadgeToUser(userID, badgeID int64) error {
    var existingID sql.NullInt64
    err := s.db.QueryRow(
        "SELECT id FROM user_badges WHERE user_id = $1 AND badge_id = $2",
        userID, badgeID,
    ).Scan(&existingID)

    if err == nil && existingID.Valid {
        PrintWarning(fmt.Sprintf("User ID %d already has badge ID %d", userID, badgeID))
        return nil
    }

    if err != nil && err != sql.ErrNoRows {
        PrintError(fmt.Sprintf("Database error: %v", err))
        return err
    }

    _, err = s.db.Exec(
        "INSERT INTO user_badges (user_id, badge_id, earned_at) VALUES ($1, $2, $3)",
        userID, badgeID, time.Now(),
    )

    if err != nil {
        PrintError(fmt.Sprintf("Failed to award badge: %v", err))
        return err
    }

    totalBadges := s.getUserBadgeCountDB(userID)
    _, err = s.db.Exec("UPDATE users SET total_badges = $1 WHERE id = $2", totalBadges, userID)

    if err != nil {
        PrintError(fmt.Sprintf("Failed to update user total_badges: %v", err))
        return err
    }

    PrintSuccess(fmt.Sprintf("Badge ID %d awarded to user ID %d", badgeID, userID))
    return nil
}

// AssignMultipleBadges assign multiple badges to multiple users
// Format: map[userID][]badgeIDs
func (s *UserBadgeSeeder) AssignMultipleBadges(assignments map[int64][]int64) error {
    tx, err := s.db.Beginx()
    if err != nil {
        PrintError(fmt.Sprintf("Failed to start transaction: %v", err))
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    totalAssigned := 0

    for userID, badgeIDs := range assignments {
        var exists bool
        err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
        if err != nil || !exists {
            PrintError(fmt.Sprintf("User ID %d not found", userID))
            continue
        }

        for _, badgeID := range badgeIDs {
            var existingID sql.NullInt64
            err := tx.QueryRow(
                "SELECT id FROM user_badges WHERE user_id = $1 AND badge_id = $2",
                userID, badgeID,
            ).Scan(&existingID)

            if err == nil && existingID.Valid {
                continue
            }

            _, err = tx.Exec(
                "INSERT INTO user_badges (user_id, badge_id, earned_at) VALUES ($1, $2, $3)",
                userID, badgeID, time.Now(),
            )

            if err != nil {
                PrintError(fmt.Sprintf("Failed to award badge: %v", err))
                return err
            }

            totalAssigned++
        }

        totalBadges := s.getUserBadgeCountTx(tx, userID)
        _, err = tx.Exec("UPDATE users SET total_badges = $1 WHERE id = $2", totalBadges, userID)

        if err != nil {
            PrintError(fmt.Sprintf("Failed to update user total_badges: %v", err))
            return err
        }
    }

    if err := tx.Commit(); err != nil {
        PrintError(fmt.Sprintf("Failed to commit transaction: %v", err))
        return err
    }

    PrintSuccess(fmt.Sprintf("Total %d badges assigned", totalAssigned))
    return nil
}