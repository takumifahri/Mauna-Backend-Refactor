package soal

import (
	"net/http"
	"strconv"

	"REFACTORING_MAUNA/internal/domain"
)

func soalIDFromPath(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, domain.NewBusinessError("INVALID_SOAL_ID", "invalid soal id", domain.ErrInvalidRequest)
	}
	return id, nil
}
