package levels

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) DeleteSubLevel(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := subLevelIDFromPath(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	if err := h.subLevelService.HardDeleteSubLevel(ctx, id); err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Sub-level deleted permanently", nil)
}
