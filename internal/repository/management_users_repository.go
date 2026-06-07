package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/dto/admin"
	"REFACTORING_MAUNA/pkg/database"
	"github.com/lib/pq"
)

type managementUsersRepository struct {
	db *database.DB
}

func NewManagementUsersRepository(db *database.DB) domain.ManagementUsersRepository {
	return &managementUsersRepository{db: db}
}

func (r *managementUsersRepository) List(ctx context.Context, filter admin.ManagementUsersFilter) (admin.ManagementUsersListResponse, error) {
	filter = normalizeManagementUsersFilter(filter)
	where, args := buildManagementUsersWhere(filter, 0)
	orderBy := managementUsersSortColumn(filter.SortBy)
	sortOrder := "DESC"
	if strings.EqualFold(filter.SortOrder, "asc") {
		sortOrder = "ASC"
	}

	query := fmt.Sprintf(`SELECT id::text, username, email, COALESCE(nama, '') AS name, COALESCE(telpon, '') AS phone,
		role::text, is_active, is_verified, total_badges, total_xp, total_quizzes_completed, total_points,
		created_at, updated_at, last_login, deleted_at
		FROM users %s ORDER BY %s %s LIMIT $%d OFFSET $%d`, where, orderBy, sortOrder, len(args)+1, len(args)+2)
	args = append(args, filter.Limit, filter.Offset)

	users := []admin.ManagementUserResponse{}
	if err := r.db.SelectContext(ctx, &users, query, args...); err != nil {
		return admin.ManagementUsersListResponse{}, domain.NewInternalError(err)
	}

	countWhere, countArgs := buildManagementUsersWhere(filter, 0)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", countWhere)
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return admin.ManagementUsersListResponse{}, domain.NewInternalError(err)
	}

	return admin.ManagementUsersListResponse{
		Users: users,
		Pagination: admin.PaginationResponse{
			Limit:  filter.Limit,
			Offset: filter.Offset,
			Total:  total,
		},
	}, nil
}

func (r *managementUsersRepository) GetByID(ctx context.Context, id string, includeDeleted bool) (admin.ManagementUserResponse, error) {
	query := `SELECT id::text, username, email, COALESCE(nama, '') AS name, COALESCE(telpon, '') AS phone,
		role::text, is_active, is_verified, total_badges, total_xp, total_quizzes_completed, total_points,
		created_at, updated_at, last_login, deleted_at
		FROM users WHERE id = $1`
	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}

	var user admin.ManagementUserResponse
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return admin.ManagementUserResponse{}, domain.ErrUserNotFound
		}
		return admin.ManagementUserResponse{}, domain.NewInternalError(err)
	}
	return user, nil
}

func (r *managementUsersRepository) Create(ctx context.Context, req admin.CreateManagementUserRequest) (string, error) {
	query := `INSERT INTO users (username, email, password, nama, telpon, role, is_active, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), $6::user_role, $7, $8, NOW(), NOW())
		RETURNING id::text`

	var id string
	if err := r.db.QueryRowContext(ctx, query, req.Username, req.Email, req.Password, req.Name, req.Phone, req.Role, boolValue(req.IsActive, true), boolValue(req.IsVerified, false)).Scan(&id); err != nil {
		return "", mapManagementUsersSQLError(err)
	}
	return id, nil
}

func (r *managementUsersRepository) Update(ctx context.Context, id string, req admin.UpdateManagementUserRequest) error {
	assignments := []string{}
	args := []any{}
	add := func(column string, value any) {
		args = append(args, value)
		assignments = append(assignments, fmt.Sprintf("%s = $%d", column, len(args)))
	}

	if req.Username != nil {
		add("username", *req.Username)
	}
	if req.Email != nil {
		add("email", *req.Email)
	}
	if req.Password != nil {
		add("password", *req.Password)
	}
	if req.Name != nil {
		add("nama", *req.Name)
	}
	if req.Phone != nil {
		phone := strings.TrimSpace(*req.Phone)
		add("telpon", sql.NullString{String: phone, Valid: phone != ""})
	}
	if req.Role != nil {
		args = append(args, *req.Role)
		assignments = append(assignments, fmt.Sprintf("role = $%d::user_role", len(args)))
	}
	if req.IsActive != nil {
		add("is_active", *req.IsActive)
	}
	if req.IsVerified != nil {
		add("is_verified", *req.IsVerified)
	}
	if len(assignments) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE users SET %s, updated_at = NOW() WHERE id = $%d AND deleted_at IS NULL", strings.Join(assignments, ", "), len(args))
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return mapManagementUsersSQLError(err)
	}
	return ensureRowsAffected(result)
}

func (r *managementUsersRepository) SoftDelete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE users SET deleted_at = NOW(), is_active = FALSE, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureRowsAffected(result)
}

func (r *managementUsersRepository) HardDelete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureRowsAffected(result)
}

func (r *managementUsersRepository) Restore(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE users SET deleted_at = NULL, is_active = TRUE, updated_at = NOW() WHERE id = $1 AND deleted_at IS NOT NULL`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureRowsAffected(result)
}

func normalizeManagementUsersFilter(filter admin.ManagementUsersFilter) admin.ManagementUsersFilter {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	return filter
}

func buildManagementUsersWhere(filter admin.ManagementUsersFilter, startIndex int) (string, []any) {
	conditions := []string{}
	args := []any{}
	next := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", startIndex+len(args))
	}

	if !filter.IncludeDeleted {
		conditions = append(conditions, "deleted_at IS NULL")
	}
	if strings.TrimSpace(filter.Query) != "" {
		value := "%" + strings.ToLower(strings.TrimSpace(filter.Query)) + "%"
		placeholder := next(value)
		conditions = append(conditions, fmt.Sprintf("(LOWER(username) LIKE %s OR LOWER(email) LIKE %s OR LOWER(COALESCE(nama, '')) LIKE %s)", placeholder, placeholder, placeholder))
	}
	if strings.TrimSpace(filter.Role) != "" {
		conditions = append(conditions, "role::text = "+next(strings.TrimSpace(filter.Role)))
	}
	if filter.IsActive != nil {
		conditions = append(conditions, "is_active = "+next(*filter.IsActive))
	}
	if filter.IsVerified != nil {
		conditions = append(conditions, "is_verified = "+next(*filter.IsVerified))
	}
	if len(conditions) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(conditions, " AND "), args
}

func managementUsersSortColumn(sortBy string) string {
	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "username":
		return "username"
	case "email":
		return "email"
	case "role":
		return "role"
	case "last_login":
		return "last_login"
	case "total_xp":
		return "total_xp"
	case "total_points":
		return "total_points"
	default:
		return "created_at"
	}
}

func mapManagementUsersSQLError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return domain.ErrUserAlreadyExists
		case "22P02":
			return domain.ErrInvalidRequest
		}
	}
	return domain.NewInternalError(err)
}

func ensureRowsAffected(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.NewInternalError(err)
	}
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func boolValue(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}
