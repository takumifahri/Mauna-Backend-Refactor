package dictionary

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) UpdateDictionary(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	id, err := dictionaryIDFromPath(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	req, err := updateDictionaryRequestFromHTTP(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.dictService.UpdateDictionary(ctx, id, req)
	if err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Dictionary updated successfully", resp)
}
