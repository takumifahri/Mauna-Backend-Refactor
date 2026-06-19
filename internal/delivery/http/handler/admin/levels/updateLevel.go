package levels

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) UpdateLevel(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	id, err := levelIDFromPath(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	req, err := updateLevelRequestFromHTTP(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.levelService.UpdateLevel(ctx, id, req)
	if err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Level updated successfully", resp)
}
