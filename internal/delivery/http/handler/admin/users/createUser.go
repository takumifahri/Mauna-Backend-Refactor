package users

import (
	"context"
	"net/http"
	"os"
	"time"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if !h.RequireAdmin(w, r) {
		return
	}

	req, uploadedAvatarPath, err := createUserRequestFromHTTP(r)
	if err != nil {
		h.WriteError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.userService.CreateUser(ctx, req)
	if err != nil {
		if uploadedAvatarPath != "" {
			_ = os.Remove(uploadedAvatarPath)
		}
		h.WriteError(w, err)
		return
	}
	h.WriteJSON(w, http.StatusCreated, "User created successfully", resp)
}
