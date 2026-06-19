package progress

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) SoftDeleteProgress(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := progressIDFromPath(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	if err := h.progressService.SoftDeleteProgress(ctx, id); err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Progress soft-deleted successfully", nil)
}
