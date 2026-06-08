package users

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"REFACTORING_MAUNA/internal/domain"
	admindto "REFACTORING_MAUNA/internal/dto/admin"
	"github.com/google/uuid"
)

const (
	createUserAvatarFormField = "avatar"
	createUserAvatarMaxSize   = 2 * 1024 * 1024
	createUserAvatarUploadDir = "uploads/avatars"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	req, uploadedAvatarPath, err := createUserRequestFromHTTP(r)
	if err != nil {
		writeAdminError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.userService.CreateUser(ctx, req)
	if err != nil {
		if uploadedAvatarPath != "" {
			_ = os.Remove(uploadedAvatarPath)
		}
		writeAdminError(w, err)
		return
	}
	writeAdminJSON(w, http.StatusCreated, "User created successfully", resp)
}

func createUserRequestFromHTTP(r *http.Request) (admindto.CreateManagementUserRequest, string, error) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		mediaType = ""
	}

	switch {
	case strings.EqualFold(mediaType, "multipart/form-data"):
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			return admindto.CreateManagementUserRequest{}, "", err
		}
		req, err := createUserRequestFromValues(r.MultipartForm.Value)
		if err != nil {
			return admindto.CreateManagementUserRequest{}, "", err
		}
		avatar, avatarURL, uploadedPath, err := saveCreateUserAvatar(r)
		if err != nil {
			return admindto.CreateManagementUserRequest{}, "", err
		}
		req.Avatar = avatar
		req.AvatarURL = avatarURL
		return req, uploadedPath, nil
	case strings.EqualFold(mediaType, "application/x-www-form-urlencoded"):
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

func createUserRequestFromValues(values map[string][]string) (admindto.CreateManagementUserRequest, error) {
	req := admindto.CreateManagementUserRequest{}

	if value, ok := firstFormValue(values, "username"); ok {
		req.Username = value
	}
	if value, ok := firstFormValue(values, "email"); ok {
		req.Email = value
	}
	if value, ok := firstFormValue(values, "password"); ok {
		req.Password = value
	}
	if value, ok := firstFormValue(values, "name"); ok {
		req.Name = value
	}
	if value, ok := firstFormValue(values, "phone"); ok {
		req.Phone = value
	}
	if value, ok := firstFormValue(values, "role"); ok {
		req.Role = value
	}
	if value, ok := firstFormValue(values, "is_active"); ok {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return admindto.CreateManagementUserRequest{}, err
		}
		req.IsActive = &parsed
	}
	if value, ok := firstFormValue(values, "is_verified"); ok {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return admindto.CreateManagementUserRequest{}, err
		}
		req.IsVerified = &parsed
	}

	return req, nil
}

func saveCreateUserAvatar(r *http.Request) (string, string, string, error) {
	file, header, err := createUserAvatarFileFromRequest(r)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return "", "", "", nil
		}
		return "", "", "", err
	}
	defer file.Close()

	if header.Size > createUserAvatarMaxSize {
		return "", "", "", domain.NewBusinessError("AVATAR_TOO_LARGE", "avatar must be 2MB or smaller", domain.ErrInvalidRequest)
	}

	extension, err := createUserAvatarExtension(file)
	if err != nil {
		return "", "", "", err
	}

	if err := os.MkdirAll(createUserAvatarUploadDir, 0755); err != nil {
		return "", "", "", err
	}

	filename := uuid.NewString() + extension
	relativePath := filepath.ToSlash(filepath.Join(createUserAvatarUploadDir, filename))
	destination, err := os.Create(relativePath)
	if err != nil {
		return "", "", "", err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, file); err != nil {
		return "", "", "", err
	}

	return filename, "/" + relativePath, relativePath, nil
}

func createUserAvatarExtension(file createUserMultipartFile) (string, error) {
	buffer := make([]byte, 512)
	size, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
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

func createUserAvatarFileFromRequest(r *http.Request) (createUserMultipartFileWithClose, *multipart.FileHeader, error) {
	file, header, err := r.FormFile(createUserAvatarFormField)
	if err == nil {
		return file, header, nil
	}
	return r.FormFile("file")
}

type createUserMultipartFile interface {
	io.Reader
	io.Seeker
}

type createUserMultipartFileWithClose interface {
	createUserMultipartFile
	io.Closer
}
