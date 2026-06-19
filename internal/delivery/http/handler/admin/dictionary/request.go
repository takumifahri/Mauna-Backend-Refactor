package dictionary

import (
	"encoding/json"
	"net/http"

	"REFACTORING_MAUNA/internal/domain"
	admin "REFACTORING_MAUNA/internal/dto/admin"
)

func createDictionaryRequestFromHTTP(r *http.Request) (admin.CreateManagementDictionaryRequest, error) {
	var req admin.CreateManagementDictionaryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return admin.CreateManagementDictionaryRequest{}, domain.ErrInvalidRequest
	}
	return req, nil
}

func updateDictionaryRequestFromHTTP(r *http.Request) (admin.UpdateManagementDictionaryRequest, error) {
	var req admin.UpdateManagementDictionaryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return admin.UpdateManagementDictionaryRequest{}, domain.ErrInvalidRequest
	}
	return req, nil
}
