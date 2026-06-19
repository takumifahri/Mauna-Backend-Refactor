package service

import (
	"context"
	"strings"

	"REFACTORING_MAUNA/internal/domain"
	admin "REFACTORING_MAUNA/internal/dto/admin"
	"REFACTORING_MAUNA/internal/usecase"
)

type managementSoalService struct {
	soalRepo domain.ManagementSoalRepository
}

var _ usecase.ManagementSoalUsecase = (*managementSoalService)(nil)

func NewManagementSoalService(soalRepo domain.ManagementSoalRepository) usecase.ManagementSoalUsecase {
	return &managementSoalService{soalRepo: soalRepo}
}

func (s *managementSoalService) ListSoal(ctx context.Context, filter admin.ManagementSoalFilter) (admin.ManagementSoalListResponse, error) {
	filter.Categories = strings.ToUpper(strings.TrimSpace(filter.Categories))
	if filter.Categories != "" && !isValidSoalType(filter.Categories) {
		return admin.ManagementSoalListResponse{}, invalidSoalTypeError()
	}
	return s.soalRepo.List(ctx, filter)
}

func (s *managementSoalService) GetSoal(ctx context.Context, id int64) (admin.ManagementSoalResponse, error) {
	return s.soalRepo.GetByID(ctx, id, false)
}

func (s *managementSoalService) CreateSoal(ctx context.Context, req admin.CreateManagementSoalRequest) (admin.ManagementSoalResponse, error) {
	question := strings.TrimSpace(req.Question)
	answer := strings.TrimSpace(req.Answer)
	categories := strings.ToUpper(strings.TrimSpace(req.Categories))
	if categories == "" {
		categories = "OPEN_CAMERA"
	}

	if err := validateSoalInput(question, answer, categories, req.DictionaryID, req.SubLevelID); err != nil {
		return admin.ManagementSoalResponse{}, err
	}

	point := req.PointGamifikasi
	if point <= 0 {
		point = 10
	}

	id, err := s.soalRepo.Create(ctx, admin.CreateManagementSoalRequest{
		Question:        question,
		Answer:          answer,
		VideoURL:        strings.TrimSpace(req.VideoURL),
		ImageURL:        strings.TrimSpace(req.ImageURL),
		DictionaryID:    req.DictionaryID,
		PointGamifikasi: point,
		SubLevelID:      req.SubLevelID,
		Categories:      categories,
	})
	if err != nil {
		return admin.ManagementSoalResponse{}, err
	}
	return s.soalRepo.GetByID(ctx, id, false)
}

func (s *managementSoalService) UpdateSoal(ctx context.Context, id int64, req admin.UpdateManagementSoalRequest) (admin.ManagementSoalResponse, error) {
	if _, err := s.soalRepo.GetByID(ctx, id, false); err != nil {
		return admin.ManagementSoalResponse{}, err
	}

	if req.Question != nil {
		q := strings.TrimSpace(*req.Question)
		if q == "" {
			return admin.ManagementSoalResponse{}, domain.NewBusinessError("INVALID_QUESTION", "question cannot be empty", domain.ErrInvalidRequest)
		}
		req.Question = &q
	}
	if req.Answer != nil {
		a := strings.TrimSpace(*req.Answer)
		if a == "" {
			return admin.ManagementSoalResponse{}, domain.NewBusinessError("INVALID_ANSWER", "answer cannot be empty", domain.ErrInvalidRequest)
		}
		req.Answer = &a
	}
	if req.Categories != nil {
		cat := strings.ToUpper(strings.TrimSpace(*req.Categories))
		if !isValidSoalType(cat) {
			return admin.ManagementSoalResponse{}, invalidSoalTypeError()
		}
		req.Categories = &cat
	}
	if req.DictionaryID != nil && *req.DictionaryID <= 0 {
		return admin.ManagementSoalResponse{}, domain.NewBusinessError("INVALID_DICTIONARY_ID", "dictionary_id must be positive", domain.ErrInvalidRequest)
	}
	if req.SubLevelID != nil && *req.SubLevelID <= 0 {
		return admin.ManagementSoalResponse{}, domain.NewBusinessError("INVALID_SUBLEVEL_ID", "sublevel_id must be positive", domain.ErrInvalidRequest)
	}

	if err := s.soalRepo.Update(ctx, id, req); err != nil {
		return admin.ManagementSoalResponse{}, err
	}
	return s.soalRepo.GetByID(ctx, id, false)
}

func (s *managementSoalService) SoftDeleteSoal(ctx context.Context, id int64) error {
	return s.soalRepo.SoftDelete(ctx, id)
}

func (s *managementSoalService) HardDeleteSoal(ctx context.Context, id int64) error {
	return s.soalRepo.HardDelete(ctx, id)
}

func (s *managementSoalService) RestoreSoal(ctx context.Context, id int64) (admin.ManagementSoalResponse, error) {
	if err := s.soalRepo.Restore(ctx, id); err != nil {
		return admin.ManagementSoalResponse{}, err
	}
	return s.soalRepo.GetByID(ctx, id, false)
}

func validateSoalInput(question, answer, categories string, dictionaryID, subLevelID int64) error {
	if strings.TrimSpace(question) == "" {
		return domain.NewBusinessError("INVALID_QUESTION", "question cannot be empty", domain.ErrInvalidRequest)
	}
	if strings.TrimSpace(answer) == "" {
		return domain.NewBusinessError("INVALID_ANSWER", "answer cannot be empty", domain.ErrInvalidRequest)
	}
	if !isValidSoalType(categories) {
		return invalidSoalTypeError()
	}
	if dictionaryID <= 0 {
		return domain.NewBusinessError("INVALID_DICTIONARY_ID", "dictionary_id is required", domain.ErrInvalidRequest)
	}
	if subLevelID <= 0 {
		return domain.NewBusinessError("INVALID_SUBLEVEL_ID", "sublevel_id is required", domain.ErrInvalidRequest)
	}
	return nil
}

func isValidSoalType(t string) bool {
	switch t {
	case "TEBAK_GAMBAR", "OPEN_CAMERA", "PILIHAN_GANDA", "MATEMATIKA":
		return true
	default:
		return false
	}
}

func invalidSoalTypeError() error {
	return domain.NewBusinessError("INVALID_SOAL_TYPE", "categories must be TEBAK_GAMBAR, OPEN_CAMERA, PILIHAN_GANDA, or MATEMATIKA", domain.ErrInvalidRequest)
}
