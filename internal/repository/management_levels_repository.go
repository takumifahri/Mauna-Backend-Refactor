package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"REFACTORING_MAUNA/internal/domain"
	admin "REFACTORING_MAUNA/internal/dto/admin"
	"REFACTORING_MAUNA/pkg/database"
	"github.com/lib/pq"
)

// ─── Level ───────────────────────────────────────────────────────────────────

type managementLevelsRepository struct {
	db *database.DB
}

func NewManagementLevelsRepository(db *database.DB) domain.ManagementLevelsRepository {
	return &managementLevelsRepository{db: db}
}

func (r *managementLevelsRepository) List(ctx context.Context, filter admin.ManagementLevelsFilter) (admin.ManagementLevelsListResponse, error) {
	filter = normalizeManagementLevelsFilter(filter)
	where, args := buildManagementLevelsWhere(filter, 0)
	orderBy := managementLevelsSortColumn(filter.SortBy)
	sortOrder := "DESC"
	if strings.EqualFold(filter.SortOrder, "asc") {
		sortOrder = "ASC"
	}

	query := fmt.Sprintf(`SELECT id, name, COALESCE(description, '') AS description,
		COALESCE(tujuan, '') AS tujuan, created_at, updated_at, deleted_at
		FROM level %s ORDER BY %s %s, id %s LIMIT $%d OFFSET $%d`,
		where, orderBy, sortOrder, sortOrder, len(args)+1, len(args)+2)
	args = append(args, filter.Limit, filter.Offset)

	items := []admin.ManagementLevelResponse{}
	if err := r.db.SelectContext(ctx, &items, query, args...); err != nil {
		return admin.ManagementLevelsListResponse{}, domain.NewInternalError(err)
	}

	countWhere, countArgs := buildManagementLevelsWhere(filter, 0)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM level %s", countWhere)
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return admin.ManagementLevelsListResponse{}, domain.NewInternalError(err)
	}

	return admin.ManagementLevelsListResponse{
		Levels: items,
		Pagination: admin.PaginationResponse{
			Limit:  filter.Limit,
			Offset: filter.Offset,
			Total:  total,
		},
	}, nil
}

func (r *managementLevelsRepository) GetByID(ctx context.Context, id int64, includeDeleted bool) (admin.ManagementLevelResponse, error) {
	query := `SELECT id, name, COALESCE(description, '') AS description,
		COALESCE(tujuan, '') AS tujuan, created_at, updated_at, deleted_at
		FROM level WHERE id = $1`
	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}

	var item admin.ManagementLevelResponse
	if err := r.db.GetContext(ctx, &item, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return admin.ManagementLevelResponse{}, domain.ErrLevelNotFound
		}
		return admin.ManagementLevelResponse{}, domain.NewInternalError(err)
	}
	return item, nil
}

func (r *managementLevelsRepository) Create(ctx context.Context, req admin.CreateManagementLevelRequest) (int64, error) {
	query := `INSERT INTO level (name, description, tujuan, created_at, updated_at)
		VALUES ($1, NULLIF($2, ''), NULLIF($3, ''), NOW(), NOW())
		RETURNING id`

	var id int64
	if err := r.db.QueryRowContext(ctx, query, req.Name, req.Description, req.Tujuan).Scan(&id); err != nil {
		return 0, mapManagementLevelsSQLError(err)
	}
	return id, nil
}

func (r *managementLevelsRepository) Update(ctx context.Context, id int64, req admin.UpdateManagementLevelRequest) error {
	assignments := []string{}
	args := []any{}
	add := func(column string, value any) {
		args = append(args, value)
		assignments = append(assignments, fmt.Sprintf("%s = $%d", column, len(args)))
	}

	if req.Name != nil {
		add("name", *req.Name)
	}
	if req.Description != nil {
		v := strings.TrimSpace(*req.Description)
		add("description", sql.NullString{String: v, Valid: v != ""})
	}
	if req.Tujuan != nil {
		v := strings.TrimSpace(*req.Tujuan)
		add("tujuan", sql.NullString{String: v, Valid: v != ""})
	}
	if len(assignments) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE level SET %s, updated_at = NOW() WHERE id = $%d AND deleted_at IS NULL",
		strings.Join(assignments, ", "), len(args))
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return mapManagementLevelsSQLError(err)
	}
	return ensureLevelRowsAffected(result)
}

func (r *managementLevelsRepository) SoftDelete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `UPDATE level SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureLevelRowsAffected(result)
}

func (r *managementLevelsRepository) HardDelete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM level WHERE id = $1`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureLevelRowsAffected(result)
}

func (r *managementLevelsRepository) Restore(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `UPDATE level SET deleted_at = NULL, updated_at = NOW() WHERE id = $1 AND deleted_at IS NOT NULL`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureLevelRowsAffected(result)
}

// ─── SubLevel ─────────────────────────────────────────────────────────────────

type managementSubLevelsRepository struct {
	db *database.DB
}

func NewManagementSubLevelsRepository(db *database.DB) domain.ManagementSubLevelsRepository {
	return &managementSubLevelsRepository{db: db}
}

func (r *managementSubLevelsRepository) List(ctx context.Context, filter admin.ManagementSubLevelsFilter) (admin.ManagementSubLevelsListResponse, error) {
	filter = normalizeManagementSubLevelsFilter(filter)
	where, args := buildManagementSubLevelsWhere(filter, 0)
	orderBy := managementSubLevelsSortColumn(filter.SortBy)
	sortOrder := "DESC"
	if strings.EqualFold(filter.SortOrder, "asc") {
		sortOrder = "ASC"
	}

	query := fmt.Sprintf(`SELECT id, name, COALESCE(description, '') AS description,
		COALESCE(tujuan, '') AS tujuan, level_id, created_at, updated_at, deleted_at
		FROM sublevel %s ORDER BY %s %s, id %s LIMIT $%d OFFSET $%d`,
		where, orderBy, sortOrder, sortOrder, len(args)+1, len(args)+2)
	args = append(args, filter.Limit, filter.Offset)

	items := []admin.ManagementSubLevelResponse{}
	if err := r.db.SelectContext(ctx, &items, query, args...); err != nil {
		return admin.ManagementSubLevelsListResponse{}, domain.NewInternalError(err)
	}

	countWhere, countArgs := buildManagementSubLevelsWhere(filter, 0)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM sublevel %s", countWhere)
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return admin.ManagementSubLevelsListResponse{}, domain.NewInternalError(err)
	}

	return admin.ManagementSubLevelsListResponse{
		SubLevels: items,
		Pagination: admin.PaginationResponse{
			Limit:  filter.Limit,
			Offset: filter.Offset,
			Total:  total,
		},
	}, nil
}

func (r *managementSubLevelsRepository) GetByID(ctx context.Context, id int64, includeDeleted bool) (admin.ManagementSubLevelResponse, error) {
	query := `SELECT id, name, COALESCE(description, '') AS description,
		COALESCE(tujuan, '') AS tujuan, level_id, created_at, updated_at, deleted_at
		FROM sublevel WHERE id = $1`
	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}

	var item admin.ManagementSubLevelResponse
	if err := r.db.GetContext(ctx, &item, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return admin.ManagementSubLevelResponse{}, domain.ErrSubLevelNotFound
		}
		return admin.ManagementSubLevelResponse{}, domain.NewInternalError(err)
	}
	return item, nil
}

func (r *managementSubLevelsRepository) Create(ctx context.Context, req admin.CreateManagementSubLevelRequest) (int64, error) {
	query := `INSERT INTO sublevel (name, description, tujuan, level_id, created_at, updated_at)
		VALUES ($1, NULLIF($2, ''), NULLIF($3, ''), $4, NOW(), NOW())
		RETURNING id`

	var id int64
	if err := r.db.QueryRowContext(ctx, query, req.Name, req.Description, req.Tujuan, req.LevelID).Scan(&id); err != nil {
		return 0, mapManagementSubLevelsSQLError(err)
	}
	return id, nil
}

func (r *managementSubLevelsRepository) Update(ctx context.Context, id int64, req admin.UpdateManagementSubLevelRequest) error {
	assignments := []string{}
	args := []any{}
	add := func(column string, value any) {
		args = append(args, value)
		assignments = append(assignments, fmt.Sprintf("%s = $%d", column, len(args)))
	}

	if req.Name != nil {
		add("name", *req.Name)
	}
	if req.Description != nil {
		v := strings.TrimSpace(*req.Description)
		add("description", sql.NullString{String: v, Valid: v != ""})
	}
	if req.Tujuan != nil {
		v := strings.TrimSpace(*req.Tujuan)
		add("tujuan", sql.NullString{String: v, Valid: v != ""})
	}
	if req.LevelID != nil {
		add("level_id", *req.LevelID)
	}
	if len(assignments) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE sublevel SET %s, updated_at = NOW() WHERE id = $%d AND deleted_at IS NULL",
		strings.Join(assignments, ", "), len(args))
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return mapManagementSubLevelsSQLError(err)
	}
	return ensureSubLevelRowsAffected(result)
}

func (r *managementSubLevelsRepository) SoftDelete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `UPDATE sublevel SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureSubLevelRowsAffected(result)
}

func (r *managementSubLevelsRepository) HardDelete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM sublevel WHERE id = $1`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureSubLevelRowsAffected(result)
}

func (r *managementSubLevelsRepository) Restore(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `UPDATE sublevel SET deleted_at = NULL, updated_at = NOW() WHERE id = $1 AND deleted_at IS NOT NULL`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureSubLevelRowsAffected(result)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func normalizeManagementLevelsFilter(filter admin.ManagementLevelsFilter) admin.ManagementLevelsFilter {
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

func buildManagementLevelsWhere(filter admin.ManagementLevelsFilter, startIndex int) (string, []any) {
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
		conditions = append(conditions, fmt.Sprintf("(LOWER(name) LIKE %s OR LOWER(COALESCE(description, '')) LIKE %s)", placeholder, placeholder))
	}
	if filter.ID > 0 {
		conditions = append(conditions, "id = "+next(filter.ID))
	}
	if strings.TrimSpace(filter.Name) != "" {
		value := "%" + strings.ToLower(strings.TrimSpace(filter.Name)) + "%"
		conditions = append(conditions, "LOWER(name) LIKE "+next(value))
	}
	if len(conditions) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(conditions, " AND "), args
}

func managementLevelsSortColumn(sortBy string) string {
	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "name":
		return "name"
	case "updated_at":
		return "updated_at"
	default:
		return "created_at"
	}
}

func mapManagementLevelsSQLError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return domain.NewBusinessError("LEVEL_ALREADY_EXISTS", "level already exists", domain.ErrInvalidRequest)
		case "23503":
			return domain.NewBusinessError("LEVEL_REF_NOT_FOUND", "referenced record not found", domain.ErrInvalidRequest)
		}
	}
	return domain.NewInternalError(err)
}

func ensureLevelRowsAffected(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.NewInternalError(err)
	}
	if rowsAffected == 0 {
		return domain.ErrLevelNotFound
	}
	return nil
}

func normalizeManagementSubLevelsFilter(filter admin.ManagementSubLevelsFilter) admin.ManagementSubLevelsFilter {
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

func buildManagementSubLevelsWhere(filter admin.ManagementSubLevelsFilter, startIndex int) (string, []any) {
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
		conditions = append(conditions, fmt.Sprintf("(LOWER(name) LIKE %s OR LOWER(COALESCE(description, '')) LIKE %s)", placeholder, placeholder))
	}
	if filter.ID > 0 {
		conditions = append(conditions, "id = "+next(filter.ID))
	}
	if strings.TrimSpace(filter.Name) != "" {
		value := "%" + strings.ToLower(strings.TrimSpace(filter.Name)) + "%"
		conditions = append(conditions, "LOWER(name) LIKE "+next(value))
	}
	if filter.LevelID > 0 {
		conditions = append(conditions, "level_id = "+next(filter.LevelID))
	}
	if len(conditions) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(conditions, " AND "), args
}

func managementSubLevelsSortColumn(sortBy string) string {
	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "name":
		return "name"
	case "level_id":
		return "level_id"
	case "updated_at":
		return "updated_at"
	default:
		return "created_at"
	}
}

func mapManagementSubLevelsSQLError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return domain.NewBusinessError("SUBLEVEL_ALREADY_EXISTS", "sublevel already exists", domain.ErrInvalidRequest)
		case "23503":
			return domain.NewBusinessError("LEVEL_NOT_FOUND", "level not found", domain.ErrInvalidRequest)
		}
	}
	return domain.NewInternalError(err)
}

func ensureSubLevelRowsAffected(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.NewInternalError(err)
	}
	if rowsAffected == 0 {
		return domain.ErrSubLevelNotFound
	}
	return nil
}
