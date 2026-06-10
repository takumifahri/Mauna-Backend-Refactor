package admin

import (
	"mime"
	"net/http"
	"strconv"
	"strings"
)

const (
	MediaTypeMultipartForm = "multipart/form-data"
	MediaTypeURLEncoded    = "application/x-www-form-urlencoded"
)

func RequestMediaType(r *http.Request) string {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return ""
	}
	return strings.ToLower(mediaType)
}

func FirstFormValue(values map[string][]string, key string) (string, bool) {
	items, ok := values[key]
	if !ok || len(items) == 0 {
		return "", false
	}
	return items[0], true
}

func ParseBool(value string) (bool, error) {
	return strconv.ParseBool(value)
}
