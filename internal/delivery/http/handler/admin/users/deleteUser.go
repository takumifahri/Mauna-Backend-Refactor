package users

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := userIDFromPath(r)
	if err != nil {
		writeAdminError(w, err)
		return
	}

	var deleteErr error
	if boolQuery(r.URL.Query().Get("hard")) {
		deleteErr = h.userService.HardDeleteUser(ctx, id)
	} else {
		deleteErr = h.userService.SoftDeleteUser(ctx, id)
	}
	if deleteErr != nil {
		writeAdminError(w, deleteErr)
		return
	}
	writeAdminJSON(w, http.StatusOK, "User deleted successfully", nil)
}
