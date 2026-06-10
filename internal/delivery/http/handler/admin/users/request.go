package users

import (
	"encoding/json"
	"net/http"
	"strings"

	adminhandler "REFACTORING_MAUNA/internal/delivery/http/handler/admin"
	"REFACTORING_MAUNA/internal/domain"
	admindto "REFACTORING_MAUNA/internal/dto/admin"
)

func createUserRequestFromHTTP(r *http.Request) (admindto.CreateManagementUserRequest, string, error) {
	mediaType := adminhandler.RequestMediaType(r)

	switch {
	case strings.EqualFold(mediaType, adminhandler.MediaTypeMultipartForm):
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			return admindto.CreateManagementUserRequest{}, "", err
		}
		req, err := createUserRequestFromValues(r.MultipartForm.Value)
		if err != nil {
			return admindto.CreateManagementUserRequest{}, "", err
		}
		avatar, avatarURL, uploadedPath, err := adminhandler.SaveAvatar(r)
		if err != nil {
			return admindto.CreateManagementUserRequest{}, "", err
		}
		req.Avatar = avatar
		req.AvatarURL = avatarURL
		return req, uploadedPath, nil
	case strings.EqualFold(mediaType, adminhandler.MediaTypeURLEncoded):
		if err := r.ParseForm(); err != nil {
			return admindto.CreateManagementUserRequest{}, "", err
		}
		req, err := createUserRequestFromValues(r.PostForm)
		return req, "", err
	default:
		var req admindto.CreateManagementUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return admindto.CreateManagementUserRequest{}, "", domain.ErrInvalidRequest
		}
		return req, "", nil
	}
}

func updateUserRequestFromHTTP(r *http.Request) (admindto.UpdateManagementUserRequest, string, error) {
	mediaType := adminhandler.RequestMediaType(r)

	switch {
	case strings.EqualFold(mediaType, adminhandler.MediaTypeMultipartForm):
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			return admindto.UpdateManagementUserRequest{}, "", err
		}
		req, err := updateUserRequestFromValues(r.MultipartForm.Value)
		if err != nil {
			return admindto.UpdateManagementUserRequest{}, "", err
		}
		avatar, avatarURL, uploadedPath, err := adminhandler.SaveAvatar(r)
		if err != nil {
			return admindto.UpdateManagementUserRequest{}, "", err
		}
		if avatar != "" || avatarURL != "" {
			req.Avatar = &avatar
			req.AvatarURL = &avatarURL
		}
		return req, uploadedPath, nil
	case strings.EqualFold(mediaType, adminhandler.MediaTypeURLEncoded):
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

func createUserRequestFromValues(values map[string][]string) (admindto.CreateManagementUserRequest, error) {
	req := admindto.CreateManagementUserRequest{}

	if value, ok := adminhandler.FirstFormValue(values, "username"); ok {
		req.Username = value
	}
	if value, ok := adminhandler.FirstFormValue(values, "email"); ok {
		req.Email = value
	}
	if value, ok := adminhandler.FirstFormValue(values, "password"); ok {
		req.Password = value
	}
	if value, ok := adminhandler.FirstFormValue(values, "name"); ok {
		req.Name = value
	}
	if value, ok := adminhandler.FirstFormValue(values, "phone"); ok {
		req.Phone = value
	}
	if value, ok := adminhandler.FirstFormValue(values, "role"); ok {
		req.Role = value
	}
	if value, ok := adminhandler.FirstFormValue(values, "is_active"); ok {
		parsed, err := adminhandler.ParseBool(value)
		if err != nil {
			return admindto.CreateManagementUserRequest{}, err
		}
		req.IsActive = &parsed
	}
	if value, ok := adminhandler.FirstFormValue(values, "is_verified"); ok {
		parsed, err := adminhandler.ParseBool(value)
		if err != nil {
			return admindto.CreateManagementUserRequest{}, err
		}
		req.IsVerified = &parsed
	}

	return req, nil
}

func updateUserRequestFromValues(values map[string][]string) (admindto.UpdateManagementUserRequest, error) {
	var req admindto.UpdateManagementUserRequest

	if value, ok := adminhandler.FirstFormValue(values, "username"); ok {
		req.Username = &value
	}
	if value, ok := adminhandler.FirstFormValue(values, "name"); ok {
		req.Name = &value
	}
	if value, ok := adminhandler.FirstFormValue(values, "avatar"); ok {
		req.Avatar = &value
	}
	if value, ok := adminhandler.FirstFormValue(values, "avatar_url"); ok {
		req.AvatarURL = &value
	}
	if value, ok := adminhandler.FirstFormValue(values, "phone"); ok {
		req.Phone = &value
	}
	if value, ok := adminhandler.FirstFormValue(values, "role"); ok {
		req.Role = &value
	}
	if value, ok := adminhandler.FirstFormValue(values, "is_active"); ok {
		parsed, err := adminhandler.ParseBool(value)
		if err != nil {
			return admindto.UpdateManagementUserRequest{}, err
		}
		req.IsActive = &parsed
	}
	if value, ok := adminhandler.FirstFormValue(values, "is_verified"); ok {
		parsed, err := adminhandler.ParseBool(value)
		if err != nil {
			return admindto.UpdateManagementUserRequest{}, err
		}
		req.IsVerified = &parsed
	}

	return req, nil
}
