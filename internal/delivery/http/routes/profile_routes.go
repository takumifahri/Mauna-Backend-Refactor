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
	authn := middleware.Auth(tokenManager)
	limit := rateLimitMiddleware.Limiter

	mux.Handle("GET /api/profile", middleware.Chain(http.HandlerFunc(profileHandler.GetProfile), authn))
	mux.Handle("PATCH /api/profile", middleware.Chain(http.HandlerFunc(profileHandler.UpdateProfile), limit(usecase.RateLimitPolicy{Limit: 10, Window: time.Minute}), authn))
	mux.Handle("POST /api/profile/avatar", middleware.Chain(http.HandlerFunc(profileHandler.UploadAvatar), limit(usecase.RateLimitPolicy{Limit: 10, Window: time.Minute}), authn))
	mux.Handle("PATCH /api/profile/password", middleware.Chain(http.HandlerFunc(profileHandler.ChangePassword), limit(usecase.RateLimitPolicy{Limit: 5, Window: time.Minute}), authn))
	mux.Handle("DELETE /api/profile", middleware.Chain(http.HandlerFunc(profileHandler.DeactivateAccount), limit(usecase.RateLimitPolicy{Limit: 3, Window: time.Minute}), authn))
}
