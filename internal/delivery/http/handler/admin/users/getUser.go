package users

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := userIDFromPath(r)
	if err != nil {
		writeAdminError(w, err)
		return
	}

	resp, err := h.userService.GetUser(ctx, id)
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeAdminJSON(w, http.StatusOK, "User retrieved successfully", resp)
}
