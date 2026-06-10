package filehandler

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestSaveImageStoresUploadedFile(t *testing.T) {
	uploadDir := t.TempDir()
	req := imageUploadRequest(t, "icon", "badge.png", tinyPNG())

	upload, err := SaveImage(req, ImageConfig{
		FormFields:  []string{"icon"},
		UploadDir:   uploadDir,
		ErrorPrefix: "BADGE_ICON",
		DisplayName: "badge icon",
		Required:    true,
	})
	if err != nil {
		t.Fatalf("SaveImage returned error: %v", err)
	}
	if upload.Filename == "" {
		t.Fatal("expected generated filename")
	}
	if !strings.HasSuffix(upload.URL, "/"+upload.Filename) {
		t.Fatalf("unexpected URL: %s", upload.URL)
	}
	if _, err := os.Stat(upload.Path); err != nil {
		t.Fatalf("expected uploaded file to exist: %v", err)
	}
}

func TestSaveImageReturnsEmptyForOptionalMissingFile(t *testing.T) {
	req := emptyMultipartRequest(t)

	upload, err := SaveImage(req, ImageConfig{FormFields: []string{"icon"}})
	if err != nil {
		t.Fatalf("SaveImage returned error: %v", err)
	}
	if upload != (ImageUpload{}) {
		t.Fatalf("expected empty upload, got %+v", upload)
	}
}

func TestSaveImageRejectsRequiredMissingFile(t *testing.T) {
	req := emptyMultipartRequest(t)

	if _, err := SaveImage(req, ImageConfig{
		FormFields:  []string{"avatar"},
		ErrorPrefix: "AVATAR",
		DisplayName: "avatar",
		Required:    true,
	}); err == nil {
		t.Fatal("expected error for missing required file")
	}
}

func TestSaveImageRejectsUnsupportedImageType(t *testing.T) {
	req := imageUploadRequest(t, "icon", "badge.txt", []byte("not an image"))

	if _, err := SaveImage(req, ImageConfig{
		FormFields:  []string{"icon"},
		UploadDir:   t.TempDir(),
		ErrorPrefix: "BADGE_ICON",
		DisplayName: "badge icon",
		Required:    true,
	}); err == nil {
		t.Fatal("expected error for unsupported image type")
	}
}

func imageUploadRequest(t *testing.T, field, filename string, content []byte) *http.Request {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(field, filename)
	if err != nil {
		t.Fatalf("CreateFormFile returned error: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("part.Write returned error: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(2 << 20); err != nil {
		t.Fatalf("ParseMultipartForm returned error: %v", err)
	}
	return req
}

func emptyMultipartRequest(t *testing.T) *http.Request {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(2 << 20); err != nil {
		t.Fatalf("ParseMultipartForm returned error: %v", err)
	}
	return req
}

func tinyPNG() []byte {
	return []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
		0x89,
	}
}
