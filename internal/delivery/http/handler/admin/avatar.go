package admin

import (
	"net/http"

	"REFACTORING_MAUNA/pkg/filehandler"
)

const (
	AvatarFormField = "avatar"
	AvatarMaxSize   = 2 * 1024 * 1024
	AvatarUploadDir = "uploads/avatars"

	BadgeIconFormField = "icon"
	BadgeIconMaxSize   = 2 * 1024 * 1024
	BadgeIconUploadDir = "uploads/badges"
)

func SaveAvatar(r *http.Request) (string, string, string, error) {
	upload, err := filehandler.SaveImage(r, filehandler.ImageConfig{
		FormFields:  []string{AvatarFormField, "file", "avatar_url"},
		UploadDir:   AvatarUploadDir,
		MaxSize:     AvatarMaxSize,
		ErrorPrefix: "AVATAR",
		DisplayName: "avatar",
	})
	return upload.Filename, upload.URL, upload.Path, err
}

func SaveBadgeIcon(r *http.Request) (string, string, string, error) {
	upload, err := filehandler.SaveImage(r, filehandler.ImageConfig{
		FormFields:  []string{BadgeIconFormField, "file", "badge_icon"},
		UploadDir:   BadgeIconUploadDir,
		MaxSize:     BadgeIconMaxSize,
		ErrorPrefix: "BADGE_ICON",
		DisplayName: "badge icon",
	})
	return upload.Filename, upload.URL, upload.Path, err
}
