package levels

import (
	"net/http"
	"strconv"

	"REFACTORING_MAUNA/internal/domain"
)

func levelIDFromPath(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, domain.NewBusinessError("INVALID_LEVEL_ID", "invalid level id", domain.ErrInvalidRequest)
	}
	return id, nil
}

func subLevelIDFromPath(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, domain.NewBusinessError("INVALID_SUBLEVEL_ID", "invalid sublevel id", domain.ErrInvalidRequest)
	}
	return id, nil
}
