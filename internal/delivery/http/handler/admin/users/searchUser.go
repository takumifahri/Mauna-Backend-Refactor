package users

import (
	"context"
	"net/http"
	"strings"
	"time"

	admindto "REFACTORING_MAUNA/internal/dto/admin"
)

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	filter := admindto.ManagementUsersFilter{
		Query:          strings.TrimSpace(r.URL.Query().Get("q")),
		Role:           strings.TrimSpace(r.URL.Query().Get("role")),
		IsActive:       optionalBool(r.URL.Query().Get("is_active")),
		IsVerified:     optionalBool(r.URL.Query().Get("is_verified")),
		IncludeDeleted: boolQuery(r.URL.Query().Get("include_deleted")),
		Limit:          intQuery(r.URL.Query().Get("limit"), 20),
		Offset:         intQuery(r.URL.Query().Get("offset"), 0),
		SortBy:         r.URL.Query().Get("sort_by"),
		SortOrder:      r.URL.Query().Get("sort_order"),
	}

	resp, err := h.userService.ListUsers(ctx, filter)
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeAdminJSON(w, http.StatusOK, "Users retrieved successfully", resp)
}
