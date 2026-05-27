package routes

import (
	"net/http"
	"time"

	profilehandler "REFACTORING_MAUNA/internal/delivery/http/handler/user/profile"
	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/usecase"
)

func RegisterProfileRoutes(mux *http.ServeMux, profileService usecase.ProfileUsecase, tokenManager usecase.TokenManager, rateLimitMiddleware *middleware.RateLimitMiddleware) {
	profileHandler := profilehandler.NewProfileHandler(profileService)

	mux.Handle("GET /api/profile", middleware.JWTAuth(tokenManager, http.HandlerFunc(profileHandler.GetProfile)))
	mux.Handle("PATCH /api/profile", rateLimitMiddleware.Limit(usecase.RateLimitPolicy{Limit: 10, Window: time.Minute}, middleware.JWTAuth(tokenManager, http.HandlerFunc(profileHandler.UpdateProfile))))
	mux.Handle("PUT /api/profile/password", rateLimitMiddleware.Limit(usecase.RateLimitPolicy{Limit: 5, Window: time.Minute}, middleware.JWTAuth(tokenManager, http.HandlerFunc(profileHandler.ChangePassword))))
	mux.Handle("DELETE /api/profile", rateLimitMiddleware.Limit(usecase.RateLimitPolicy{Limit: 3, Window: time.Minute}, middleware.JWTAuth(tokenManager, http.HandlerFunc(profileHandler.DeactivateAccount))))
}
