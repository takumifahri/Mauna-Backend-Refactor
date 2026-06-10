package profile

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/dto"
)

const (
	avatarFormField = "avatar"
	avatarMaxSize   = 2 * 1024 * 1024
	avatarBodyLimit = avatarMaxSize + 1024*1024
	avatarUploadDir = "uploads/avatars"
)

func (h *Handler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeProfileErrorResponse(w, http.StatusUnauthorized, domain.ErrUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, avatarBodyLimit)
	if err := r.ParseMultipartForm(avatarMaxSize); err != nil {
		writeProfileErrorResponse(w, http.StatusBadRequest, domain.NewBusinessError("INVALID_AVATAR", "avatar must be sent as multipart/form-data and 2MB or smaller", domain.ErrInvalidRequest))
		return
	}

	filename, avatarURL, cleanup, err := saveAvatarFromRequest(r)
	if err != nil {
		writeProfileErrorResponse(w, domain.ErrorToStatusCode(err), err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.profileService.UpdateAvatar(ctx, claims.UserID, filename, avatarURL)
	if err != nil {
		cleanup()
		writeProfileErrorResponse(w, domain.ErrorToStatusCode(err), err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.Response{
		Status:    "success",
		Message:   "Avatar uploaded successfully",
		Data:      resp,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
