package levels

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) CreateLevel(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	req, err := createLevelRequestFromHTTP(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.levelService.CreateLevel(ctx, req)
	if err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusCreated, "Level created successfully", resp)
}
