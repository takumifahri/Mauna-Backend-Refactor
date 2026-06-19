package soal

import (
	"encoding/json"
	"net/http"

	"REFACTORING_MAUNA/internal/domain"
	admin "REFACTORING_MAUNA/internal/dto/admin"
)

func createSoalRequestFromHTTP(r *http.Request) (admin.CreateManagementSoalRequest, error) {
	var req admin.CreateManagementSoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return admin.CreateManagementSoalRequest{}, domain.ErrInvalidRequest
	}
	return req, nil
}

func updateSoalRequestFromHTTP(r *http.Request) (admin.UpdateManagementSoalRequest, error) {
	var req admin.UpdateManagementSoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return admin.UpdateManagementSoalRequest{}, domain.ErrInvalidRequest
	}
	return req, nil
}
