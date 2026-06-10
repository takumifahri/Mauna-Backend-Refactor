package badges

import (
	"context"
	"net/http"
	"os"
	"time"
)

func (h *Handler) CreateBadge(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	req, uploadedIconPath, err := createBadgeRequestFromHTTP(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.badgeService.CreateBadge(ctx, req)
	if err != nil {
		if uploadedIconPath != "" {
			_ = os.Remove(uploadedIconPath)
		}
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusCreated, "Badge created successfully", resp)
}
