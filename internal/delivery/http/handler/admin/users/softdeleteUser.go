package users

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) SoftDeleteUser(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := userIDFromPath(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	err = h.userService.SoftDeleteUser(ctx, id)
	if err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "User soft-deleted successfully", nil)
}
