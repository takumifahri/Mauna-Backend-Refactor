package filehandler

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"REFACTORING_MAUNA/internal/domain"
	"github.com/google/uuid"
)

type ImageConfig struct {
	FormFields  []string
	UploadDir   string
	MaxSize     int64
	ErrorPrefix string
	DisplayName string
	Required    bool
}

type ImageUpload struct {
	Filename string
	URL      string
	Path     string
}

func SaveImage(r *http.Request, cfg ImageConfig) (ImageUpload, error) {
	cfg = normalizeImageConfig(cfg)

	file, header, err := imageFileFromRequest(r, cfg.FormFields)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) && !cfg.Required {
			return ImageUpload{}, nil
		}
		if errors.Is(err, http.ErrMissingFile) {
			return ImageUpload{}, domain.NewBusinessError("INVALID_"+cfg.ErrorPrefix, cfg.DisplayName+" file is required", domain.ErrInvalidRequest)
		}
		return ImageUpload{}, err
	}
	defer file.Close()

	if header.Size > cfg.MaxSize {
		return ImageUpload{}, domain.NewBusinessError(cfg.ErrorPrefix+"_TOO_LARGE", cfg.DisplayName+" must be 2MB or smaller", domain.ErrInvalidRequest)
	}

	extension, err := imageExtension(file, cfg)
	if err != nil {
		return ImageUpload{}, err
	}

	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		return ImageUpload{}, domain.NewInternalError(err)
	}

	filename := uuid.NewString() + extension
	relativePath := filepath.ToSlash(filepath.Join(cfg.UploadDir, filename))
	destination, err := os.Create(relativePath)
	if err != nil {
		return ImageUpload{}, domain.NewInternalError(err)
	}
	defer destination.Close()

	if _, err := io.Copy(destination, file); err != nil {
		return ImageUpload{}, domain.NewInternalError(err)
	}

	return ImageUpload{
		Filename: filename,
		URL:      "/" + relativePath,
		Path:     relativePath,
	}, nil
}

func HasMultipartFile(r *http.Request, fields ...string) bool {
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return false
	}
	for _, field := range fields {
		if files := r.MultipartForm.File[field]; len(files) > 0 {
			return true
		}
	}
	return false
}

func normalizeImageConfig(cfg ImageConfig) ImageConfig {
	if cfg.MaxSize == 0 {
		cfg.MaxSize = 2 * 1024 * 1024
	}
	if cfg.ErrorPrefix == "" {
		cfg.ErrorPrefix = "IMAGE"
	}
	cfg.ErrorPrefix = strings.ToUpper(strings.TrimSpace(cfg.ErrorPrefix))
	if cfg.DisplayName == "" {
		cfg.DisplayName = strings.ToLower(strings.ReplaceAll(cfg.ErrorPrefix, "_", " "))
	}
	if len(cfg.FormFields) == 0 {
		cfg.FormFields = []string{"file"}
	}
	return cfg
}

func imageExtension(file multipartFile, cfg ImageConfig) (string, error) {
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
		return "", domain.NewBusinessError("INVALID_"+cfg.ErrorPrefix+"_TYPE", cfg.DisplayName+" must be a jpeg, png, or webp image", domain.ErrInvalidRequest)
	}
}

func imageFileFromRequest(r *http.Request, fields []string) (multipartFileWithClose, *multipart.FileHeader, error) {
	var lastErr error
	for _, field := range fields {
		file, header, err := r.FormFile(field)
		if err == nil {
			return file, header, nil
		}
		lastErr = err
	}
	return nil, nil, lastErr
}

type multipartFile interface {
	io.Reader
	io.Seeker
}

type multipartFileWithClose interface {
	multipartFile
	io.Closer
}
