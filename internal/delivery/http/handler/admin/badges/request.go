package badges

import (
	"encoding/json"
	"net/http"
	"strings"

	adminhandler "REFACTORING_MAUNA/internal/delivery/http/handler/admin"
	"REFACTORING_MAUNA/internal/domain"
	admindto "REFACTORING_MAUNA/internal/dto/admin"
)

func createBadgeRequestFromHTTP(r *http.Request) (admindto.CreateManagementBadgeRequest, string, error) {
	mediaType := adminhandler.RequestMediaType(r)

	switch {
	case strings.EqualFold(mediaType, adminhandler.MediaTypeMultipartForm):
		if err := r.ParseMultipartForm(2 << 20); err != nil {
			return admindto.CreateManagementBadgeRequest{}, "", err
		}
		req := createBadgeRequestFromValues(r.MultipartForm.Value)
		_, iconURL, uploadedPath, err := adminhandler.SaveBadgeIcon(r)
		if err != nil {
			return admindto.CreateManagementBadgeRequest{}, "", err
		}
		if iconURL != "" {
			req.Icon = iconURL
		}
		return req, uploadedPath, nil
	case strings.EqualFold(mediaType, adminhandler.MediaTypeURLEncoded):
		if err := r.ParseForm(); err != nil {
			return admindto.CreateManagementBadgeRequest{}, "", err
		}
		return createBadgeRequestFromValues(r.PostForm), "", nil
	default:
		var req admindto.CreateManagementBadgeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return admindto.CreateManagementBadgeRequest{}, "", domain.ErrInvalidRequest
		}
		return req, "", nil
	}
}

func updateBadgeRequestFromHTTP(r *http.Request) (admindto.UpdateManagementBadgeRequest, string, error) {
	mediaType := adminhandler.RequestMediaType(r)

	switch {
	case strings.EqualFold(mediaType, adminhandler.MediaTypeMultipartForm):
		if err := r.ParseMultipartForm(2 << 20); err != nil {
			return admindto.UpdateManagementBadgeRequest{}, "", err
		}
		req := updateBadgeRequestFromValues(r.MultipartForm.Value)
		_, iconURL, uploadedPath, err := adminhandler.SaveBadgeIcon(r)
		if err != nil {
			return admindto.UpdateManagementBadgeRequest{}, "", err
		}
		if iconURL != "" {
			req.Icon = &iconURL
		}
		return req, uploadedPath, nil
	case strings.EqualFold(mediaType, adminhandler.MediaTypeURLEncoded):
		if err := r.ParseForm(); err != nil {
			return admindto.UpdateManagementBadgeRequest{}, "", err
		}
		return updateBadgeRequestFromValues(r.PostForm), "", nil
	default:
		var req admindto.UpdateManagementBadgeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return admindto.UpdateManagementBadgeRequest{}, "", domain.ErrInvalidRequest
		}
		return req, "", nil
	}
}

func createBadgeRequestFromValues(values map[string][]string) admindto.CreateManagementBadgeRequest {
	req := admindto.CreateManagementBadgeRequest{}
	if value, ok := adminhandler.FirstFormValue(values, "name"); ok {
		req.Name = value
	}
	if value, ok := adminhandler.FirstFormValue(values, "nama"); ok {
		req.Name = value
	}
	if value, ok := adminhandler.FirstFormValue(values, "description"); ok {
		req.Description = value
	}
	if value, ok := adminhandler.FirstFormValue(values, "deskripsi"); ok {
		req.Description = value
	}
	if value, ok := adminhandler.FirstFormValue(values, "icon"); ok {
		req.Icon = value
	}
	if value, ok := adminhandler.FirstFormValue(values, "level"); ok {
		req.Level = value
	}
	return req
}

func updateBadgeRequestFromValues(values map[string][]string) admindto.UpdateManagementBadgeRequest {
	req := admindto.UpdateManagementBadgeRequest{}
	if value, ok := adminhandler.FirstFormValue(values, "name"); ok {
		req.Name = &value
	}
	if value, ok := adminhandler.FirstFormValue(values, "nama"); ok {
		req.Name = &value
	}
	if value, ok := adminhandler.FirstFormValue(values, "description"); ok {
		req.Description = &value
	}
	if value, ok := adminhandler.FirstFormValue(values, "deskripsi"); ok {
		req.Description = &value
	}
	if value, ok := adminhandler.FirstFormValue(values, "icon"); ok {
		req.Icon = &value
	}
	if value, ok := adminhandler.FirstFormValue(values, "level"); ok {
		req.Level = &value
	}
	return req
}
