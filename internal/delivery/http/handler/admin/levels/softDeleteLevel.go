package levels

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) SoftDeleteLevel(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := levelIDFromPath(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	if err := h.levelService.SoftDeleteLevel(ctx, id); err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Level soft-deleted successfully", nil)
}
