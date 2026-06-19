package soal

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) DeleteSoal(w http.ResponseWriter, r *http.Request) {
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

	if err := h.soalService.HardDeleteSoal(ctx, id); err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Soal deleted permanently", nil)
}
