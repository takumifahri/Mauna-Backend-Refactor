package levels

import (
	"context"
	"net/http"
	"strings"
	"time"

	admin "REFACTORING_MAUNA/internal/dto/admin"
)

func (h *Handler) ListSubLevels(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	filter := admin.ManagementSubLevelsFilter{
		Query:          strings.TrimSpace(r.URL.Query().Get("q")),
		ID:             int64Query(r.URL.Query().Get("id")),
		Name:           strings.TrimSpace(r.URL.Query().Get("name")),
		LevelID:        int64Query(r.URL.Query().Get("level_id")),
		IncludeDeleted: boolQuery(r.URL.Query().Get("include_deleted")),
		Limit:          intQuery(r.URL.Query().Get("limit"), 20),
		Offset:         intQuery(r.URL.Query().Get("offset"), 0),
		SortBy:         r.URL.Query().Get("sort_by"),
		SortOrder:      r.URL.Query().Get("sort_order"),
	}

	resp, err := h.subLevelService.ListSubLevels(ctx, filter)
	if err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Sub-levels retrieved successfully", resp)
}
