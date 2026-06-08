package profile

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/dto"
	"github.com/google/uuid"
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

	file, header, err := avatarFileFromRequest(r)
	if err != nil {
		writeProfileErrorResponse(w, http.StatusBadRequest, domain.NewBusinessError("INVALID_AVATAR", "avatar file is required", domain.ErrInvalidRequest))
		return
	}
	defer file.Close()

	if header.Size > avatarMaxSize {
		writeProfileErrorResponse(w, http.StatusBadRequest, domain.NewBusinessError("AVATAR_TOO_LARGE", "avatar must be 2MB or smaller", domain.ErrInvalidRequest))
		return
	}

	extension, err := avatarExtension(file)
	if err != nil {
		writeProfileErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := os.MkdirAll(avatarUploadDir, 0755); err != nil {
		writeProfileErrorResponse(w, http.StatusInternalServerError, domain.NewInternalError(err))
		return
	}

	filename := uuid.NewString() + extension
	relativePath := filepath.ToSlash(filepath.Join(avatarUploadDir, filename))
	destination, err := os.Create(relativePath)
	if err != nil {
		writeProfileErrorResponse(w, http.StatusInternalServerError, domain.NewInternalError(err))
		return
	}
	defer destination.Close()

	if _, err := io.Copy(destination, file); err != nil {
		writeProfileErrorResponse(w, http.StatusInternalServerError, domain.NewInternalError(err))
		return
	}

	avatarURL := "/" + relativePath

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.profileService.UpdateAvatar(ctx, claims.UserID, filename, avatarURL)
	if err != nil {
		_ = os.Remove(relativePath)
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

func avatarExtension(file multipartFile) (string, error) {
	buffer := make([]byte, 512)
	size, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", domain.NewInternalError(err)
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", domain.NewInternalError(err)
	}

	switch strings.ToLower(http.DetectContentType(buffer[:size])) {
	case "image/jpeg":
		return ".jpg", nil
	case "image/png":
		return ".png", nil
	case "image/webp":
		return ".webp", nil
	default:
		return "", domain.NewBusinessError("INVALID_AVATAR_TYPE", "avatar must be a jpeg, png, or webp image", domain.ErrInvalidRequest)
	}
}

type multipartFile interface {
	io.Reader
	io.Seeker
}

func avatarFileFromRequest(r *http.Request) (multipartFileWithClose, *multipart.FileHeader, error) {
	file, header, err := r.FormFile(avatarFormField)
	if err == nil {
		return file, header, nil
	}
	return r.FormFile("file")
}

type multipartFileWithClose interface {
	multipartFile
	io.Closer
}
