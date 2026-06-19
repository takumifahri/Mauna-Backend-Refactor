package dictionary

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) SoftDeleteDictionary(w http.ResponseWriter, r *http.Request) {
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

	if err := h.dictService.SoftDeleteDictionary(ctx, id); err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Dictionary soft-deleted successfully", nil)
}
