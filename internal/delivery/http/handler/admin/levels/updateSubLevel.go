package levels

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) UpdateSubLevel(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	id, err := subLevelIDFromPath(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	req, err := updateSubLevelRequestFromHTTP(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.subLevelService.UpdateSubLevel(ctx, id, req)
	if err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Sub-level updated successfully", resp)
}
