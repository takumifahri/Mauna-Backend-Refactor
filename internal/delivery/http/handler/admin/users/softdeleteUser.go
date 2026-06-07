package users

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) RestoreUser(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.userService.RestoreUser(ctx, r.PathValue("id"))
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeAdminJSON(w, http.StatusOK, "User restored successfully", resp)
}
