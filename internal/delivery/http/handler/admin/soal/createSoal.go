package soal

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) CreateSoal(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	req, err := createSoalRequestFromHTTP(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.soalService.CreateSoal(ctx, req)
	if err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusCreated, "Soal created successfully", resp)
}
