package badges

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) SoftDeleteBadge(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := badgeIDFromPath(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	if err := h.badgeService.SoftDeleteBadge(ctx, id); err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Badge soft-deleted successfully", nil)
}
