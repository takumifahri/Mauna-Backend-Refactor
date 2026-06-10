package profile

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/dto"
	"REFACTORING_MAUNA/pkg/filehandler"
)

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeProfileErrorResponse(w, http.StatusUnauthorized, domain.ErrUnauthorized)
		return
	}

	req, hasAvatarUpload, err := updateProfileRequestFromHTTP(w, r)
	if err != nil {
		slog.WarnContext(r.Context(), "profile_update_decode_failed", slog.Any("error", err))
		writeProfileErrorResponse(w, http.StatusBadRequest, domain.ErrInvalidRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.profileService.UpdateProfile(ctx, claims.UserID, req)
	if err != nil {
		statusCode := domain.ErrorToStatusCode(err)
		writeProfileErrorResponse(w, statusCode, err)
		return
	}

	if hasAvatarUpload {
		filename, avatarURL, cleanup, err := saveAvatarFromRequest(r)
		if err != nil {
			writeProfileErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		resp, err = h.profileService.UpdateAvatar(ctx, claims.UserID, filename, avatarURL)
		if err != nil {
			cleanup()
			statusCode := domain.ErrorToStatusCode(err)
			writeProfileErrorResponse(w, statusCode, err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.Response{
		Status:    "success",
		Message:   "Profile updated successfully",
		Data:      resp,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func updateProfileRequestFromHTTP(w http.ResponseWriter, r *http.Request) (dto.UpdateProfileRequest, bool, error) {
	contentType := strings.ToLower(r.Header.Get("Content-Type"))

	if strings.HasPrefix(contentType, "multipart/form-data") {
		r.Body = http.MaxBytesReader(w, r.Body, avatarBodyLimit)
		if err := r.ParseMultipartForm(avatarMaxSize); err != nil {
			return dto.UpdateProfileRequest{}, false, err
		}
		return updateProfileRequestFromForm(r), hasAvatarFile(r), nil
	}

	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			return dto.UpdateProfileRequest{}, false, err
		}
		return updateProfileRequestFromForm(r), false, nil
	}

	var req dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return dto.UpdateProfileRequest{}, false, err
	}
	return req, false, nil
}

func updateProfileRequestFromForm(r *http.Request) dto.UpdateProfileRequest {
	return dto.UpdateProfileRequest{
		Username:  firstProfileFormValue(r, "username"),
		Name:      firstProfileFormValue(r, "name", "nama"),
		Phone:     firstProfileFormValue(r, "phone", "telpon"),
		AvatarURL: firstProfileFormValue(r, "avatar_url"),
		Bio:       firstProfileFormValue(r, "bio"),
	}
}

func firstProfileFormValue(r *http.Request, keys ...string) *string {
	for _, key := range keys {
		values, ok := r.Form[key]
		if !ok || len(values) == 0 {
			continue
		}
		value := values[0]
		return &value
	}
	return nil
}

func hasAvatarFile(r *http.Request) bool {
	return filehandler.HasMultipartFile(r, avatarFormField, "file", "avatar_url")
}

func saveAvatarFromRequest(r *http.Request) (string, string, func(), error) {
	upload, err := filehandler.SaveImage(r, filehandler.ImageConfig{
		FormFields:  []string{avatarFormField, "file", "avatar_url"},
		UploadDir:   avatarUploadDir,
		MaxSize:     avatarMaxSize,
		ErrorPrefix: "AVATAR",
		DisplayName: "avatar",
		Required:    true,
	})
	if err != nil {
		return "", "", func() {}, err
	}

	return upload.Filename, upload.URL, func() { _ = os.Remove(upload.Path) }, nil
}
