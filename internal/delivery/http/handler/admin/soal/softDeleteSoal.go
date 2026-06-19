package soal

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) SoftDeleteSoal(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := soalIDFromPath(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	if err := h.soalService.SoftDeleteSoal(ctx, id); err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Soal soft-deleted successfully", nil)
}
