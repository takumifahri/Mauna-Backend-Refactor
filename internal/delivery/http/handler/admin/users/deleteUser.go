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

	var err error
	if boolQuery(r.URL.Query().Get("hard")) {
		err = h.userService.HardDeleteUser(ctx, r.PathValue("id"))
	} else {
		err = h.userService.SoftDeleteUser(ctx, r.PathValue("id"))
	}
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeAdminJSON(w, http.StatusOK, "User deleted successfully", nil)
}
