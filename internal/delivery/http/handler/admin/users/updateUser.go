package users

import (
	"context"
	"net/http"
	"os"
	"time"
)

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	req, uploadedAvatarPath, err := updateUserRequestFromHTTP(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := userIDFromPath(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	resp, err := h.userService.UpdateUser(ctx, id, req)
	if err != nil {
		if uploadedAvatarPath != "" {
			_ = os.Remove(uploadedAvatarPath)
		}
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusOK, "User updated successfully", resp)
}
