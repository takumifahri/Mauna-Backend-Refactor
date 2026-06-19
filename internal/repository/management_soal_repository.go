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

type managementSoalRepository struct {
	db *database.DB
}

func NewManagementSoalRepository(db *database.DB) domain.ManagementSoalRepository {
	return &managementSoalRepository{db: db}
}

func (r *managementSoalRepository) List(ctx context.Context, filter admin.ManagementSoalFilter) (admin.ManagementSoalListResponse, error) {
	filter = normalizeManagementSoalFilter(filter)
	where, args := buildManagementSoalWhere(filter, 0)
	orderBy := managementSoalSortColumn(filter.SortBy)
	sortOrder := "DESC"
	if strings.EqualFold(filter.SortOrder, "asc") {
		sortOrder = "ASC"
	}

	query := fmt.Sprintf(`SELECT id, question, answer,
		COALESCE(video_url, '') AS video_url, COALESCE(image_url, '') AS image_url,
		dictionary_id, point_gamifikasi, sublevel_id, categories::text AS categories,
		created_at, updated_at, deleted_at
		FROM soal %s ORDER BY %s %s, id %s LIMIT $%d OFFSET $%d`,
		where, orderBy, sortOrder, sortOrder, len(args)+1, len(args)+2)
	args = append(args, filter.Limit, filter.Offset)

	items := []admin.ManagementSoalResponse{}
	if err := r.db.SelectContext(ctx, &items, query, args...); err != nil {
		return admin.ManagementSoalListResponse{}, domain.NewInternalError(err)
	}

	countWhere, countArgs := buildManagementSoalWhere(filter, 0)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM soal %s", countWhere)
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return admin.ManagementSoalListResponse{}, domain.NewInternalError(err)
	}

	return admin.ManagementSoalListResponse{
		Soal: items,
		Pagination: admin.PaginationResponse{
			Limit:  filter.Limit,
			Offset: filter.Offset,
			Total:  total,
		},
	}, nil
}

func (r *managementSoalRepository) GetByID(ctx context.Context, id int64, includeDeleted bool) (admin.ManagementSoalResponse, error) {
	query := `SELECT id, question, answer,
		COALESCE(video_url, '') AS video_url, COALESCE(image_url, '') AS image_url,
		dictionary_id, point_gamifikasi, sublevel_id, categories::text AS categories,
		created_at, updated_at, deleted_at
		FROM soal WHERE id = $1`
	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}

	var item admin.ManagementSoalResponse
	if err := r.db.GetContext(ctx, &item, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return admin.ManagementSoalResponse{}, domain.ErrQuestionNotFound
		}
		return admin.ManagementSoalResponse{}, domain.NewInternalError(err)
	}
	return item, nil
}

func (r *managementSoalRepository) Create(ctx context.Context, req admin.CreateManagementSoalRequest) (int64, error) {
	point := req.PointGamifikasi
	if point <= 0 {
		point = 10
	}

	query := `INSERT INTO soal (question, answer, video_url, image_url, dictionary_id, point_gamifikasi, sublevel_id, categories, created_at, updated_at)
		VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''), $5, $6, $7, $8::soal_type, NOW(), NOW())
		RETURNING id`

	var id int64
	if err := r.db.QueryRowContext(ctx, query,
		req.Question, req.Answer, req.VideoURL, req.ImageURL,
		req.DictionaryID, point, req.SubLevelID, req.Categories,
	).Scan(&id); err != nil {
		return 0, mapManagementSoalSQLError(err)
	}
	return id, nil
}

func (r *managementSoalRepository) Update(ctx context.Context, id int64, req admin.UpdateManagementSoalRequest) error {
	assignments := []string{}
	args := []any{}
	add := func(column string, value any) {
		args = append(args, value)
		assignments = append(assignments, fmt.Sprintf("%s = $%d", column, len(args)))
	}

	if req.Question != nil {
		add("question", *req.Question)
	}
	if req.Answer != nil {
		add("answer", *req.Answer)
	}
	if req.VideoURL != nil {
		v := strings.TrimSpace(*req.VideoURL)
		add("video_url", sql.NullString{String: v, Valid: v != ""})
	}
	if req.ImageURL != nil {
		v := strings.TrimSpace(*req.ImageURL)
		add("image_url", sql.NullString{String: v, Valid: v != ""})
	}
	if req.DictionaryID != nil {
		add("dictionary_id", *req.DictionaryID)
	}
	if req.PointGamifikasi != nil {
		add("point_gamifikasi", *req.PointGamifikasi)
	}
	if req.SubLevelID != nil {
		add("sublevel_id", *req.SubLevelID)
	}
	if req.Categories != nil {
		args = append(args, *req.Categories)
		assignments = append(assignments, fmt.Sprintf("categories = $%d::soal_type", len(args)))
	}
	if len(assignments) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE soal SET %s, updated_at = NOW() WHERE id = $%d AND deleted_at IS NULL",
		strings.Join(assignments, ", "), len(args))
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return mapManagementSoalSQLError(err)
	}
	return ensureSoalRowsAffected(result)
}

func (r *managementSoalRepository) SoftDelete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `UPDATE soal SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureSoalRowsAffected(result)
}

func (r *managementSoalRepository) HardDelete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM soal WHERE id = $1`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureSoalRowsAffected(result)
}

func (r *managementSoalRepository) Restore(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `UPDATE soal SET deleted_at = NULL, updated_at = NOW() WHERE id = $1 AND deleted_at IS NOT NULL`, id)
	if err != nil {
		return domain.NewInternalError(err)
	}
	return ensureSoalRowsAffected(result)
}

func normalizeManagementSoalFilter(filter admin.ManagementSoalFilter) admin.ManagementSoalFilter {
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

func buildManagementSoalWhere(filter admin.ManagementSoalFilter, startIndex int) (string, []any) {
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
		conditions = append(conditions, fmt.Sprintf("(LOWER(question) LIKE %s OR LOWER(answer) LIKE %s)", placeholder, placeholder))
	}
	if filter.ID > 0 {
		conditions = append(conditions, "id = "+next(filter.ID))
	}
	if filter.SubLevelID > 0 {
		conditions = append(conditions, "sublevel_id = "+next(filter.SubLevelID))
	}
	if filter.DictionaryID > 0 {
		conditions = append(conditions, "dictionary_id = "+next(filter.DictionaryID))
	}
	if strings.TrimSpace(filter.Categories) != "" {
		conditions = append(conditions, "categories::text = "+next(strings.ToUpper(strings.TrimSpace(filter.Categories))))
	}
	if len(conditions) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(conditions, " AND "), args
}

func managementSoalSortColumn(sortBy string) string {
	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "sublevel_id":
		return "sublevel_id"
	case "dictionary_id":
		return "dictionary_id"
	case "categories":
		return "categories"
	case "point_gamifikasi":
		return "point_gamifikasi"
	case "updated_at":
		return "updated_at"
	default:
		return "created_at"
	}
}

func mapManagementSoalSQLError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23503":
			return domain.NewBusinessError("REF_NOT_FOUND", "referenced dictionary or sublevel not found", domain.ErrInvalidRequest)
		case "22P02":
			return domain.NewBusinessError("INVALID_SOAL_TYPE", "invalid question category", domain.ErrInvalidRequest)
		}
	}
	return domain.NewInternalError(err)
}

func ensureSoalRowsAffected(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.NewInternalError(err)
	}
	if rowsAffected == 0 {
		return domain.ErrQuestionNotFound
	}
	return nil
}
