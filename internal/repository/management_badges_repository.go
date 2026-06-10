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

type managementBadgesRepository struct {
	db *database.DB
}

func NewManagementBadgesRepository(db *database.DB) domain.ManagementBadgesRepository {
	return &managementBadgesRepository{db: db}
}

func (r *managementBadgesRepository) List(ctx context.Context, filter admin.ManagementBadgesFilter) (admin.ManagementBadgesListResponse, error) {
	filter = normalizeManagementBadgesFilter(filter)
	where, args := buildManagementBadgesWhere(filter, 0)
	orderBy := managementBadgesSortColumn(filter.SortBy)
	sortOrder := "DESC"
	if strings.EqualFold(filter.SortOrder, "asc") {
		sortOrder = "ASC"
	}

	query := fmt.Sprintf(`SELECT id, nama AS name, COALESCE(deskripsi, '') AS description,
		COALESCE(icon, '') AS icon, level::text AS level, created_at, updated_at, deleted_at
		FROM badges %s ORDER BY %s %s LIMIT $%d OFFSET $%d`, where, orderBy, sortOrder, len(args)+1, len(args)+2)
	args = append(args, filter.Limit, filter.Offset)

	badges := []admin.ManagementBadgeResponse{}
	if err := r.db.SelectContext(ctx, &badges, query, args...); err != nil {
		return admin.ManagementBadgesListResponse{}, domain.NewInternalError(err)
	}

	countWhere, countArgs := buildManagementBadgesWhere(filter, 0)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM badges %s", countWhere)
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return admin.ManagementBadgesListResponse{}, domain.NewInternalError(err)
	}

	return admin.ManagementBadgesListResponse{
		Badges: badges,
		Pagination: admin.PaginationResponse{
			Limit:  filter.Limit,
			Offset: filter.Offset,
			Total:  total,
		},
	}, nil
}

func (r *managementBadgesRepository) GetByID(ctx context.Context, id int64, includeDeleted bool) (admin.ManagementBadgeResponse, error) {
	query := `SELECT id, nama AS name, COALESCE(deskripsi, '') AS description,
		COALESCE(icon, '') AS icon, level::text AS level, created_at, updated_at, deleted_at
		FROM badges WHERE id = $1`
	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}

	var badge admin.ManagementBadgeResponse
	if err := r.db.GetContext(ctx, &badge, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return admin.ManagementBadgeResponse{}, domain.ErrBadgeNotFound
		}
		return admin.ManagementBadgeResponse{}, domain.NewInternalError(err)
	}
	return badge, nil
}

func (r *managementBadgesRepository) Create(ctx context.Context, req admin.CreateManagementBadgeRequest) (int64, error) {
	query := `INSERT INTO badges (nama, deskripsi, icon, level, created_at, updated_at)
		VALUES ($1, NULLIF($2, ''), NULLIF($3, ''), $4::difficulty_level, NOW(), NOW())
		RETURNING id`

	var id int64
	if err := r.db.QueryRowContext(ctx, query, req.Name, req.Description, req.Icon, req.Level).Scan(&id); err != nil {
		return 0, mapManagementBadgesSQLError(err)
	}
	return id, nil
}

func (r *managementBadgesRepository) Update(ctx context.Context, id int64, req admin.UpdateManagementBadgeRequest) error {
	assignments := []string{}
	args := []any{}
	add := func(column string, value any) {
		args = append(args, value)
		assignments = append(assignments, fmt.Sprintf("%s = $%d", column, len(args)))
	}

	if req.Name != nil {
		add("nama", *req.Name)
	}
	if req.Description != nil {
		description := strings.TrimSpace(*req.Description)
		add("deskripsi", sql.NullString{String: description, Valid: description != ""})
	}
	if req.Icon != nil {
		icon := strings.TrimSpace(*req.Icon)
		add("icon", sql.NullString{String: icon, Valid: icon != ""})
	}
	if req.Level != nil {
		args = append(args, *req.Level)
		assignments = append(assignments, fmt.Sprintf("level = $%d::difficulty_level", len(args)))
	}
	if len(assignments) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE badges SET %s, updated_at = NOW() WHERE id = $%d AND deleted_at IS NULL", strings.Join(assignments, ", "), len(args))
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return mapManagementBadgesSQLError(err)
	}
	return ensureBadgeRowsAffected(result)
}

func (r *managementBadgesRepository) SoftDelete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `UPDATE badges SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureBadgeRowsAffected(result)
}

func (r *managementBadgesRepository) HardDelete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM badges WHERE id = $1`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureBadgeRowsAffected(result)
}

func (r *managementBadgesRepository) Restore(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `UPDATE badges SET deleted_at = NULL, updated_at = NOW() WHERE id = $1 AND deleted_at IS NOT NULL`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureBadgeRowsAffected(result)
}

func normalizeManagementBadgesFilter(filter admin.ManagementBadgesFilter) admin.ManagementBadgesFilter {
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

func buildManagementBadgesWhere(filter admin.ManagementBadgesFilter, startIndex int) (string, []any) {
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
		conditions = append(conditions, fmt.Sprintf("(LOWER(nama) LIKE %s OR LOWER(COALESCE(deskripsi, '')) LIKE %s)", placeholder, placeholder))
	}
	if filter.ID > 0 {
		conditions = append(conditions, "id = "+next(filter.ID))
	}
	if strings.TrimSpace(filter.Name) != "" {
		value := "%" + strings.ToLower(strings.TrimSpace(filter.Name)) + "%"
		conditions = append(conditions, "LOWER(nama) LIKE "+next(value))
	}
	if strings.TrimSpace(filter.Level) != "" {
		conditions = append(conditions, "level::text = "+next(strings.ToLower(strings.TrimSpace(filter.Level))))
	}
	if len(conditions) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(conditions, " AND "), args
}

func managementBadgesSortColumn(sortBy string) string {
	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "name", "nama":
		return "nama"
	case "level":
		return "level"
	case "updated_at":
		return "updated_at"
	default:
		return "created_at"
	}
}

func mapManagementBadgesSQLError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return domain.NewBusinessError("BADGE_ALREADY_EXISTS", "badge already exists", domain.ErrInvalidRequest)
		case "22P02":
			return domain.ErrInvalidRequest
		}
	}
	return domain.NewInternalError(err)
}

func ensureBadgeRowsAffected(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.NewInternalError(err)
	}
	if rowsAffected == 0 {
		return domain.ErrBadgeNotFound
	}
	return nil
}
