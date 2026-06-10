package badges

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) DeleteBadge(w http.ResponseWriter, r *http.Request) {
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

	var deleteErr error
	if boolQuery(r.URL.Query().Get("hard")) {
		deleteErr = h.badgeService.HardDeleteBadge(ctx, id)
	} else {
		deleteErr = h.badgeService.SoftDeleteBadge(ctx, id)
	}
	if deleteErr != nil {
		h.WriteError(w, deleteErr)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Badge deleted successfully", nil)
}
