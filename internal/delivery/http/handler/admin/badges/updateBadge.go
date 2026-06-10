package badges

import (
	"context"
	"net/http"
	"os"
	"time"
)

func (h *Handler) UpdateBadge(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	id, err := badgeIDFromPath(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	req, uploadedIconPath, err := updateBadgeRequestFromHTTP(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.badgeService.UpdateBadge(ctx, id, req)
	if err != nil {
		if uploadedIconPath != "" {
			_ = os.Remove(uploadedIconPath)
		}
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Badge updated successfully", resp)
}
