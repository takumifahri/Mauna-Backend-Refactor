package auth

import (
	"net/http"
	"os"
	"strings"
	"time"

	"REFACTORING_MAUNA/pkg/security"
)

const (
	accessTokenCookieName  = "access_token"
	refreshTokenCookieName = "refresh_token"
)

func setAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	setTokenCookie(w, accessTokenCookieName, accessToken, security.AccessTokenDuration)
	setTokenCookie(w, refreshTokenCookieName, refreshToken, security.RefreshTokenDuration)
}

func clearAuthCookies(w http.ResponseWriter) {
	clearTokenCookie(w, accessTokenCookieName)
	clearTokenCookie(w, refreshTokenCookieName)
}

func getRefreshToken(r *http.Request, bodyToken string) string {
	if bodyToken != "" {
		return bodyToken
	}

	cookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func setTokenCookie(w http.ResponseWriter, name, value string, maxAge time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   int(maxAge.Seconds()),
		Expires:  time.Now().Add(maxAge),
		HttpOnly: true,
		Secure:   cookieSecure(),
		SameSite: http.SameSiteLaxMode,
	})
}

func clearTokenCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   cookieSecure(),
		SameSite: http.SameSiteLaxMode,
	})
}

func cookieSecure() bool {
	value := strings.ToLower(os.Getenv("COOKIE_SECURE"))
	return value == "true" || value == "1" || value == "yes"
}
