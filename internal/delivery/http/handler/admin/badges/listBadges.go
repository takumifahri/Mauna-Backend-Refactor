package badges

import (
	"context"
	"net/http"
	"strings"
	"time"

	admindto "REFACTORING_MAUNA/internal/dto/admin"
)

func (h *Handler) ListBadges(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	filter := admindto.ManagementBadgesFilter{
		Query:          strings.TrimSpace(r.URL.Query().Get("q")),
		ID:             int64Query(r.URL.Query().Get("id")),
		Name:           strings.TrimSpace(r.URL.Query().Get("name")),
		Level:          strings.TrimSpace(r.URL.Query().Get("level")),
		IncludeDeleted: boolQuery(r.URL.Query().Get("include_deleted")),
		Limit:          intQuery(r.URL.Query().Get("limit"), 20),
		Offset:         intQuery(r.URL.Query().Get("offset"), 0),
		SortBy:         r.URL.Query().Get("sort_by"),
		SortOrder:      r.URL.Query().Get("sort_order"),
	}

	resp, err := h.badgeService.ListBadges(ctx, filter)
	if err != nil {
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "Badges retrieved successfully", resp)
}
