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

type managementDictionaryRepository struct {
	db *database.DB
}

func NewManagementDictionaryRepository(db *database.DB) domain.ManagementDictionaryRepository {
	return &managementDictionaryRepository{db: db}
}

func (r *managementDictionaryRepository) List(ctx context.Context, filter admin.ManagementDictionaryFilter) (admin.ManagementDictionaryListResponse, error) {
	filter = normalizeManagementDictionaryFilter(filter)
	where, args := buildManagementDictionaryWhere(filter, 0)
	orderBy := managementDictionarySortColumn(filter.SortBy)
	sortOrder := "DESC"
	if strings.EqualFold(filter.SortOrder, "asc") {
		sortOrder = "ASC"
	}

	query := fmt.Sprintf(`SELECT id, word_text, definition, category::text AS category,
		COALESCE(image_url_ref, '') AS image_url_ref, COALESCE(video_url, '') AS video_url,
		created_at, updated_at, deleted_at
		FROM kamus %s ORDER BY %s %s, id %s LIMIT $%d OFFSET $%d`,
		where, orderBy, sortOrder, sortOrder, len(args)+1, len(args)+2)
	args = append(args, filter.Limit, filter.Offset)

	items := []admin.ManagementDictionaryResponse{}
	if err := r.db.SelectContext(ctx, &items, query, args...); err != nil {
		return admin.ManagementDictionaryListResponse{}, domain.NewInternalError(err)
	}

	countWhere, countArgs := buildManagementDictionaryWhere(filter, 0)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM kamus %s", countWhere)
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return admin.ManagementDictionaryListResponse{}, domain.NewInternalError(err)
	}

	return admin.ManagementDictionaryListResponse{
		Dictionaries: items,
		Pagination: admin.PaginationResponse{
			Limit:  filter.Limit,
			Offset: filter.Offset,
			Total:  total,
		},
	}, nil
}

func (r *managementDictionaryRepository) GetByID(ctx context.Context, id int64, includeDeleted bool) (admin.ManagementDictionaryResponse, error) {
	query := `SELECT id, word_text, definition, category::text AS category,
		COALESCE(image_url_ref, '') AS image_url_ref, COALESCE(video_url, '') AS video_url,
		created_at, updated_at, deleted_at
		FROM kamus WHERE id = $1`
	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}

	var item admin.ManagementDictionaryResponse
	if err := r.db.GetContext(ctx, &item, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return admin.ManagementDictionaryResponse{}, domain.ErrDictionaryNotFound
		}
		return admin.ManagementDictionaryResponse{}, domain.NewInternalError(err)
	}
	return item, nil
}

func (r *managementDictionaryRepository) Create(ctx context.Context, req admin.CreateManagementDictionaryRequest) (int64, error) {
	query := `INSERT INTO kamus (word_text, definition, category, image_url_ref, video_url, created_at, updated_at)
		VALUES ($1, $2, $3::kamus_category, NULLIF($4, ''), NULLIF($5, ''), NOW(), NOW())
		RETURNING id`

	var id int64
	if err := r.db.QueryRowContext(ctx, query, req.WordText, req.Definition, req.Category, req.ImageURLRef, req.VideoURL).Scan(&id); err != nil {
		return 0, mapManagementDictionarySQLError(err)
	}
	return id, nil
}

func (r *managementDictionaryRepository) Update(ctx context.Context, id int64, req admin.UpdateManagementDictionaryRequest) error {
	assignments := []string{}
	args := []any{}
	add := func(column string, value any) {
		args = append(args, value)
		assignments = append(assignments, fmt.Sprintf("%s = $%d", column, len(args)))
	}

	if req.WordText != nil {
		add("word_text", *req.WordText)
	}
	if req.Definition != nil {
		add("definition", *req.Definition)
	}
	if req.Category != nil {
		args = append(args, *req.Category)
		assignments = append(assignments, fmt.Sprintf("category = $%d::kamus_category", len(args)))
	}
	if req.ImageURLRef != nil {
		v := strings.TrimSpace(*req.ImageURLRef)
		add("image_url_ref", sql.NullString{String: v, Valid: v != ""})
	}
	if req.VideoURL != nil {
		v := strings.TrimSpace(*req.VideoURL)
		add("video_url", sql.NullString{String: v, Valid: v != ""})
	}
	if len(assignments) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE kamus SET %s, updated_at = NOW() WHERE id = $%d AND deleted_at IS NULL",
		strings.Join(assignments, ", "), len(args))
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return mapManagementDictionarySQLError(err)
	}
	return ensureDictionaryRowsAffected(result)
}

func (r *managementDictionaryRepository) SoftDelete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `UPDATE kamus SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureDictionaryRowsAffected(result)
}

func (r *managementDictionaryRepository) HardDelete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM kamus WHERE id = $1`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureDictionaryRowsAffected(result)
}

func (r *managementDictionaryRepository) Restore(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `UPDATE kamus SET deleted_at = NULL, updated_at = NOW() WHERE id = $1 AND deleted_at IS NOT NULL`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureDictionaryRowsAffected(result)
}

func normalizeManagementDictionaryFilter(filter admin.ManagementDictionaryFilter) admin.ManagementDictionaryFilter {
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

func buildManagementDictionaryWhere(filter admin.ManagementDictionaryFilter, startIndex int) (string, []any) {
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
		conditions = append(conditions, fmt.Sprintf("(LOWER(word_text) LIKE %s OR LOWER(definition) LIKE %s)", placeholder, placeholder))
	}
	if filter.ID > 0 {
		conditions = append(conditions, "id = "+next(filter.ID))
	}
	if strings.TrimSpace(filter.WordText) != "" {
		value := "%" + strings.ToLower(strings.TrimSpace(filter.WordText)) + "%"
		conditions = append(conditions, "LOWER(word_text) LIKE "+next(value))
	}
	if strings.TrimSpace(filter.Category) != "" {
		conditions = append(conditions, "category::text = "+next(strings.ToUpper(strings.TrimSpace(filter.Category))))
	}
	if len(conditions) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(conditions, " AND "), args
}

func managementDictionarySortColumn(sortBy string) string {
	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "word_text":
		return "word_text"
	case "category":
		return "category"
	case "updated_at":
		return "updated_at"
	default:
		return "created_at"
	}
}

func mapManagementDictionarySQLError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return domain.NewBusinessError("DICTIONARY_ALREADY_EXISTS", "word already exists in dictionary", domain.ErrInvalidRequest)
		case "22P02":
			return domain.ErrInvalidRequest
		}
	}
	return domain.NewInternalError(err)
}

func ensureDictionaryRowsAffected(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.NewInternalError(err)
	}
	if rowsAffected == 0 {
		return domain.ErrDictionaryNotFound
	}
	return nil
}
