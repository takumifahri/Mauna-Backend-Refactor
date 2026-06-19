package progress

import (
	"context"
	"net/http"
	"strings"
	"time"

	admin "REFACTORING_MAUNA/internal/dto/admin"
)

func (h *Handler) ListProgress(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	filter := admin.ManagementProgressFilter{
		ID:             int64Query(r.URL.Query().Get("id")),
		UserID:         strings.TrimSpace(r.URL.Query().Get("user_id")),
		SubLevelID:     int64Query(r.URL.Query().Get("sublevel_id")),
		Status:         strings.TrimSpace(r.URL.Query().Get("status")),
		IncludeDeleted: boolQuery(r.URL.Query().Get("include_deleted")),
		Limit:          intQuery(r.URL.Query().Get("limit"), 20),
		Offset:         intQuery(r.URL.Query().Get("offset"), 0),
		SortBy:         r.URL.Query().Get("sort_by"),
		SortOrder:      r.URL.Query().Get("sort_order"),
	}

	resp, err := h.progressService.ListProgress(ctx, filter)
	if err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Progress retrieved successfully", resp)
}
