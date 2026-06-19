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
)

type managementProgressRepository struct {
	db *database.DB
}

func NewManagementProgressRepository(db *database.DB) domain.ManagementProgressRepository {
	return &managementProgressRepository{db: db}
}

func (r *managementProgressRepository) List(ctx context.Context, filter admin.ManagementProgressFilter) (admin.ManagementProgressListResponse, error) {
	filter = normalizeManagementProgressFilter(filter)
	where, args := buildManagementProgressWhere(filter, 0)
	orderBy := managementProgressSortColumn(filter.SortBy)
	sortOrder := "DESC"
	if strings.EqualFold(filter.SortOrder, "asc") {
		sortOrder = "ASC"
	}

	query := fmt.Sprintf(`SELECT id, user_id::text, sublevel_id, status::text AS status,
		total_questions, correct_answers, score, stars, completion_percentage,
		attempts, best_score, best_stars, is_unlocked,
		created_at, updated_at, deleted_at, last_attempt, completed_at
		FROM progress %s ORDER BY %s %s, id %s LIMIT $%d OFFSET $%d`,
		where, orderBy, sortOrder, sortOrder, len(args)+1, len(args)+2)
	args = append(args, filter.Limit, filter.Offset)

	items := []admin.ManagementProgressResponse{}
	if err := r.db.SelectContext(ctx, &items, query, args...); err != nil {
		return admin.ManagementProgressListResponse{}, domain.NewInternalError(err)
	}

	countWhere, countArgs := buildManagementProgressWhere(filter, 0)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM progress %s", countWhere)
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return admin.ManagementProgressListResponse{}, domain.NewInternalError(err)
	}

	return admin.ManagementProgressListResponse{
		Progress: items,
		Pagination: admin.PaginationResponse{
			Limit:  filter.Limit,
			Offset: filter.Offset,
			Total:  total,
		},
	}, nil
}

func (r *managementProgressRepository) GetByID(ctx context.Context, id int64, includeDeleted bool) (admin.ManagementProgressResponse, error) {
	query := `SELECT id, user_id::text, sublevel_id, status::text AS status,
		total_questions, correct_answers, score, stars, completion_percentage,
		attempts, best_score, best_stars, is_unlocked,
		created_at, updated_at, deleted_at, last_attempt, completed_at
		FROM progress WHERE id = $1`
	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}

	var item admin.ManagementProgressResponse
	if err := r.db.GetContext(ctx, &item, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return admin.ManagementProgressResponse{}, domain.ErrProgressNotFound
		}
		return admin.ManagementProgressResponse{}, domain.NewInternalError(err)
	}
	return item, nil
}

func (r *managementProgressRepository) SoftDelete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `UPDATE progress SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureProgressRowsAffected(result)
}

func (r *managementProgressRepository) HardDelete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM progress WHERE id = $1`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureProgressRowsAffected(result)
}

func (r *managementProgressRepository) Restore(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `UPDATE progress SET deleted_at = NULL, updated_at = NOW() WHERE id = $1 AND deleted_at IS NOT NULL`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureProgressRowsAffected(result)
}

func normalizeManagementProgressFilter(filter admin.ManagementProgressFilter) admin.ManagementProgressFilter {
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

func buildManagementProgressWhere(filter admin.ManagementProgressFilter, startIndex int) (string, []any) {
	conditions := []string{}
	args := []any{}
	next := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", startIndex+len(args))
	}

	if !filter.IncludeDeleted {
		conditions = append(conditions, "deleted_at IS NULL")
	}
	if filter.ID > 0 {
		conditions = append(conditions, "id = "+next(filter.ID))
	}
	if strings.TrimSpace(filter.UserID) != "" {
		conditions = append(conditions, "user_id = "+next(filter.UserID)+"::uuid")
	}
	if filter.SubLevelID > 0 {
		conditions = append(conditions, "sublevel_id = "+next(filter.SubLevelID))
	}
	if strings.TrimSpace(filter.Status) != "" {
		conditions = append(conditions, "status::text = "+next(strings.ToLower(strings.TrimSpace(filter.Status))))
	}
	if len(conditions) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(conditions, " AND "), args
}

func managementProgressSortColumn(sortBy string) string {
	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "user_id":
		return "user_id"
	case "sublevel_id":
		return "sublevel_id"
	case "status":
		return "status"
	case "score":
		return "score"
	case "updated_at":
		return "updated_at"
	default:
		return "created_at"
	}
}

func ensureProgressRowsAffected(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.NewInternalError(err)
	}
	if rowsAffected == 0 {
		return domain.ErrProgressNotFound
	}
	return nil
}
