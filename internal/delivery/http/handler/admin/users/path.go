package users

import (
	"net/http"

	"REFACTORING_MAUNA/internal/domain"
	"github.com/google/uuid"
)

func userIDFromPath(r *http.Request) (string, error) {
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		return "", domain.NewBusinessError("INVALID_USER_ID", "invalid user id", domain.ErrInvalidRequest)
	}
	return id, nil
}
