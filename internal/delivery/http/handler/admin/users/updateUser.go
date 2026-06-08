package users

import (
	"context"
	"encoding/json"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"REFACTORING_MAUNA/internal/domain"
	admindto "REFACTORING_MAUNA/internal/dto/admin"
)

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	req, uploadedAvatarPath, err := updateUserRequestFromHTTP(r)
	if err != nil {
		writeAdminError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := userIDFromPath(r)
	if err != nil {
		writeAdminError(w, err)
		return
	}

	resp, err := h.userService.UpdateUser(ctx, id, req)
	if err != nil {
		if uploadedAvatarPath != "" {
			_ = os.Remove(uploadedAvatarPath)
		}
		writeAdminError(w, err)
		return
	}
	writeAdminJSON(w, http.StatusOK, "User updated successfully", resp)
}

func updateUserRequestFromHTTP(r *http.Request) (admindto.UpdateManagementUserRequest, string, error) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		mediaType = ""
	}

	switch {
	case strings.EqualFold(mediaType, "multipart/form-data"):
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			return admindto.UpdateManagementUserRequest{}, "", err
		}
		req, err := updateUserRequestFromValues(r.MultipartForm.Value)
		if err != nil {
			return admindto.UpdateManagementUserRequest{}, "", err
		}
		avatar, avatarURL, uploadedPath, err := saveCreateUserAvatar(r)
		if err != nil {
			return admindto.UpdateManagementUserRequest{}, "", err
		}
		if avatar != "" || avatarURL != "" {
			req.Avatar = &avatar
			req.AvatarURL = &avatarURL
		}
		return req, uploadedPath, nil
	case strings.EqualFold(mediaType, "application/x-www-form-urlencoded"):
		if err := r.ParseForm(); err != nil {
			return admindto.UpdateManagementUserRequest{}, "", err
		}
		req, err := updateUserRequestFromValues(r.PostForm)
		return req, "", err
	default:
		var req admindto.UpdateManagementUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return admindto.UpdateManagementUserRequest{}, "", domain.ErrInvalidRequest
		}
		return req, "", nil
	}
}

func updateUserRequestFromValues(values map[string][]string) (admindto.UpdateManagementUserRequest, error) {
	var req admindto.UpdateManagementUserRequest

	if value, ok := firstFormValue(values, "username"); ok {
		req.Username = &value
	}
	// if value, ok := firstFormValue(values, "email"); ok {
	// 	req.Email = &value
	// }
	// if value, ok := firstFormValue(values, "password"); ok {
	// 	req.Password = &value
	// }
	if value, ok := firstFormValue(values, "name"); ok {
		req.Name = &value
	}
	if value, ok := firstFormValue(values, "avatar"); ok {
		req.Avatar = &value
	}
	if value, ok := firstFormValue(values, "avatar_url"); ok {
		req.AvatarURL = &value
	}
	if value, ok := firstFormValue(values, "phone"); ok {
		req.Phone = &value
	}
	if value, ok := firstFormValue(values, "role"); ok {
		req.Role = &value
	}
	if value, ok := firstFormValue(values, "is_active"); ok {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return admindto.UpdateManagementUserRequest{}, err
		}
		req.IsActive = &parsed
	}
	if value, ok := firstFormValue(values, "is_verified"); ok {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return admindto.UpdateManagementUserRequest{}, err
		}
		req.IsVerified = &parsed
	}

	return req, nil
}

func firstFormValue(values map[string][]string, key string) (string, bool) {
	items, ok := values[key]
	if !ok || len(items) == 0 {
		return "", false
	}
	return items[0], true
}
