package soal

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) UpdateSoal(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	id, err := soalIDFromPath(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	req, err := updateSoalRequestFromHTTP(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.soalService.UpdateSoal(ctx, id, req)
	if err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Soal updated successfully", resp)
}
