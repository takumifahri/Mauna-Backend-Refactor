package soal

import (
	"context"
	"net/http"
	"strings"
	"time"

	admin "REFACTORING_MAUNA/internal/dto/admin"
)

func (h *Handler) ListSoal(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	filter := admin.ManagementSoalFilter{
		Query:          strings.TrimSpace(r.URL.Query().Get("q")),
		ID:             int64Query(r.URL.Query().Get("id")),
		SubLevelID:     int64Query(r.URL.Query().Get("sublevel_id")),
		DictionaryID:   int64Query(r.URL.Query().Get("dictionary_id")),
		Categories:     strings.TrimSpace(r.URL.Query().Get("categories")),
		IncludeDeleted: boolQuery(r.URL.Query().Get("include_deleted")),
		Limit:          intQuery(r.URL.Query().Get("limit"), 20),
		Offset:         intQuery(r.URL.Query().Get("offset"), 0),
		SortBy:         r.URL.Query().Get("sort_by"),
		SortOrder:      r.URL.Query().Get("sort_order"),
	}

	resp, err := h.soalService.ListSoal(ctx, filter)
	if err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Soal retrieved successfully", resp)
}
