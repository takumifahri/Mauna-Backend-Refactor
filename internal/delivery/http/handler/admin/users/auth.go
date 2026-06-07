package users

import (
	"net/http"

	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/domain"
)

func (h *Handler) requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeAdminError(w, domain.ErrUnauthorized)
		return false
	}
	if claims.Role != "admin" {
		writeAdminError(w, domain.ErrForbidden)
		return false
	}
	return true
}
