package dictionary

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) GetDictionary(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := dictionaryIDFromPath(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	resp, err := h.dictService.GetDictionary(ctx, id)
	if err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Dictionary retrieved successfully", resp)
}
