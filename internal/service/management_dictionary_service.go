package service

import (
	"context"
	"strings"

	"REFACTORING_MAUNA/internal/domain"
	admin "REFACTORING_MAUNA/internal/dto/admin"
	"REFACTORING_MAUNA/internal/usecase"
)

type managementDictionaryService struct {
	dictRepo domain.ManagementDictionaryRepository
}

var _ usecase.ManagementDictionaryUsecase = (*managementDictionaryService)(nil)

func NewManagementDictionaryService(dictRepo domain.ManagementDictionaryRepository) usecase.ManagementDictionaryUsecase {
	return &managementDictionaryService{dictRepo: dictRepo}
}

func (s *managementDictionaryService) ListDictionary(ctx context.Context, filter admin.ManagementDictionaryFilter) (admin.ManagementDictionaryListResponse, error) {
	filter.Category = strings.ToUpper(strings.TrimSpace(filter.Category))
	if filter.Category != "" && !isValidKamusCategory(filter.Category) {
		return admin.ManagementDictionaryListResponse{}, invalidKamusCategoryError()
	}
	return s.dictRepo.List(ctx, filter)
}

func (s *managementDictionaryService) GetDictionary(ctx context.Context, id int64) (admin.ManagementDictionaryResponse, error) {
	return s.dictRepo.GetByID(ctx, id, false)
}

func (s *managementDictionaryService) CreateDictionary(ctx context.Context, req admin.CreateManagementDictionaryRequest) (admin.ManagementDictionaryResponse, error) {
	wordText := strings.TrimSpace(req.WordText)
	definition := strings.TrimSpace(req.Definition)
	category := strings.ToUpper(strings.TrimSpace(req.Category))
	if category == "" {
		category = "ALPHABET"
	}

	if err := validateDictionaryInput(wordText, definition, category); err != nil {
		return admin.ManagementDictionaryResponse{}, err
	}

	id, err := s.dictRepo.Create(ctx, admin.CreateManagementDictionaryRequest{
		WordText:    wordText,
		Definition:  definition,
		Category:    category,
		ImageURLRef: strings.TrimSpace(req.ImageURLRef),
		VideoURL:    strings.TrimSpace(req.VideoURL),
	})
	if err != nil {
		return admin.ManagementDictionaryResponse{}, err
	}
	return s.dictRepo.GetByID(ctx, id, false)
}

func (s *managementDictionaryService) UpdateDictionary(ctx context.Context, id int64, req admin.UpdateManagementDictionaryRequest) (admin.ManagementDictionaryResponse, error) {
	if _, err := s.dictRepo.GetByID(ctx, id, false); err != nil {
		return admin.ManagementDictionaryResponse{}, err
	}

	if req.WordText != nil {
		wordText := strings.TrimSpace(*req.WordText)
		if wordText == "" {
			return admin.ManagementDictionaryResponse{}, domain.NewBusinessError("INVALID_WORD_TEXT", "word text cannot be empty", domain.ErrInvalidRequest)
		}
		req.WordText = &wordText
	}
	if req.Definition != nil {
		definition := strings.TrimSpace(*req.Definition)
		if definition == "" {
			return admin.ManagementDictionaryResponse{}, domain.NewBusinessError("INVALID_DEFINITION", "definition cannot be empty", domain.ErrInvalidRequest)
		}
		req.Definition = &definition
	}
	if req.Category != nil {
		category := strings.ToUpper(strings.TrimSpace(*req.Category))
		if !isValidKamusCategory(category) {
			return admin.ManagementDictionaryResponse{}, invalidKamusCategoryError()
		}
		req.Category = &category
	}

	if err := s.dictRepo.Update(ctx, id, req); err != nil {
		return admin.ManagementDictionaryResponse{}, err
	}
	return s.dictRepo.GetByID(ctx, id, false)
}

func (s *managementDictionaryService) SoftDeleteDictionary(ctx context.Context, id int64) error {
	return s.dictRepo.SoftDelete(ctx, id)
}

func (s *managementDictionaryService) HardDeleteDictionary(ctx context.Context, id int64) error {
	return s.dictRepo.HardDelete(ctx, id)
}

func (s *managementDictionaryService) RestoreDictionary(ctx context.Context, id int64) (admin.ManagementDictionaryResponse, error) {
	if err := s.dictRepo.Restore(ctx, id); err != nil {
		return admin.ManagementDictionaryResponse{}, err
	}
	return s.dictRepo.GetByID(ctx, id, false)
}

func validateDictionaryInput(wordText, definition, category string) error {
	if strings.TrimSpace(wordText) == "" {
		return domain.NewBusinessError("INVALID_WORD_TEXT", "word text cannot be empty", domain.ErrInvalidRequest)
	}
	if strings.TrimSpace(definition) == "" {
		return domain.NewBusinessError("INVALID_DEFINITION", "definition cannot be empty", domain.ErrInvalidRequest)
	}
	if !isValidKamusCategory(category) {
		return invalidKamusCategoryError()
	}
	return nil
}

func isValidKamusCategory(category string) bool {
	switch category {
	case "ALPHABET", "NUMBERS", "IMBUHAN", "KOSAKATA":
		return true
	default:
		return false
	}
}

func invalidKamusCategoryError() error {
	return domain.NewBusinessError("INVALID_CATEGORY", "category must be ALPHABET, NUMBERS, IMBUHAN, or KOSAKATA", domain.ErrInvalidRequest)
}
