package levels

import (
	"encoding/json"
	"net/http"

	"REFACTORING_MAUNA/internal/domain"
	admin "REFACTORING_MAUNA/internal/dto/admin"
)

func createLevelRequestFromHTTP(r *http.Request) (admin.CreateManagementLevelRequest, error) {
	var req admin.CreateManagementLevelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return admin.CreateManagementLevelRequest{}, domain.ErrInvalidRequest
	}
	return req, nil
}

func updateLevelRequestFromHTTP(r *http.Request) (admin.UpdateManagementLevelRequest, error) {
	var req admin.UpdateManagementLevelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return admin.UpdateManagementLevelRequest{}, domain.ErrInvalidRequest
	}
	return req, nil
}

func createSubLevelRequestFromHTTP(r *http.Request) (admin.CreateManagementSubLevelRequest, error) {
	var req admin.CreateManagementSubLevelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return admin.CreateManagementSubLevelRequest{}, domain.ErrInvalidRequest
	}
	return req, nil
}

func updateSubLevelRequestFromHTTP(r *http.Request) (admin.UpdateManagementSubLevelRequest, error) {
	var req admin.UpdateManagementSubLevelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return admin.UpdateManagementSubLevelRequest{}, domain.ErrInvalidRequest
	}
	return req, nil
}
