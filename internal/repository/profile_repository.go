package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"REFACTORING_MAUNA/internal/domain/constant"
	"REFACTORING_MAUNA/internal/domain/entities"
)

type ProfileRepository struct {
	DB *sql.DB
}

func (r *ProfileRepository) GetProfile(ctx context.Context, userID string) (*entities.User, error) {
	query := "SELECT id, email, nama, telpon, role, is_active, is_verified, bio FROM users WHERE id = ?"
	row := r.DB.QueryRowContext(ctx, query, userID)

	var user entities.User
	err := row.Scan(&user.ID, &user.Email, &user.Nama, &user.Telpon, &user.Role, &user.IsActive, &user.IsVerified, &user.Bio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return &user, nil
}

func (r *ProfileRepository) UpdateProfile(ctx context.Context, userID string, updatedData *entities.User) error {
	query := "UPDATE users SET email = ?, nama = ?, telpon = ?, bio = ? WHERE id = ?"
	_, err := r.DB.ExecContext(ctx, query, updatedData.Email, updatedData.Nama, updatedData.Telpon, updatedData.Bio, userID)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
	return nil
}

func (r *ProfileRepository) DeactivateAccount(ctx context.Context, userID string) error {
	query := "UPDATE users SET is_active = ? WHERE id = ?"
	_, err := r.DB.ExecContext(ctx, query, false, userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate account: %w", err)
	}
	return nil
}

func (r *ProfileRepository) ChangePassword(ctx context.Context, userID string, newPasswordHash string) error {
	query := "UPDATE users SET password = ? WHERE id = ?"
	_, err := r.DB.ExecContext(ctx, query, newPasswordHash, userID)
	if err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}
	return nil
}

func (r *ProfileRepository) GetUserRole(ctx context.Context, userID string) (constant.UserRole, error) {
	query := "SELECT role FROM users WHERE id = ?"
	row := r.DB.QueryRowContext(ctx, query, userID)

	var role constant.UserRole
	err := row.Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("failed to get user role: %w", err)
	}

	return role, nil
}

func (r *ProfileRepository) GetUserActivity(ctx context.Context, userID string) (int, error) {
	query := "SELECT total_quizzes_completed FROM users WHERE id = ?"
	row := r.DB.QueryRowContext(ctx, query, userID)

	var totalQuizzes int
	err := row.Scan(&totalQuizzes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user not found")
		}
		return 0, fmt.Errorf("failed to get user activity: %w", err)
	}

	return totalQuizzes, nil
}

func (r *ProfileRepository) GetUserBadges(ctx context.Context, userID string) (int, error) {
	query := "SELECT total_badges FROM users WHERE id = ?"
	row := r.DB.QueryRowContext(ctx, query, userID)

	var totalBadges int
	err := row.Scan(&totalBadges)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user not found")
		}
		return 0, fmt.Errorf("failed to get user badges: %w", err)
	}

	return totalBadges, nil
}

func (r *ProfileRepository) GetUserProgress(ctx context.Context, userID string) (float64, error) {
	query := "SELECT progress_percentage FROM users WHERE id = ?"
	row := r.DB.QueryRowContext(ctx, query, userID)

	var progress float64
	err := row.Scan(&progress)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user not found")
		}
		return 0, fmt.Errorf("failed to get user progress: %w", err)
	}

	return progress, nil
}
