package profile

import "REFACTORING_MAUNA/internal/usecase"

type Handler struct {
	profileService usecase.ProfileUsecase
}

func NewProfileHandler(profileService usecase.ProfileUsecase) *Handler {
	return &Handler{
		profileService: profileService,
	}
}
