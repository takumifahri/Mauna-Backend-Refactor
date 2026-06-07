package http

import (
	"fmt"
	"net/http"
	"time"

	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/delivery/http/routes"
	"REFACTORING_MAUNA/internal/repository"
	"REFACTORING_MAUNA/internal/service"
	"REFACTORING_MAUNA/pkg/database"
	"REFACTORING_MAUNA/pkg/mailer"
	"REFACTORING_MAUNA/pkg/security"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Route represents a route definition
type Route struct {
	Path    string
	Method  string
	Handler string
}

// ← RENAME: routes → availableRoutes
var availableRoutes = []Route{
	{Path: "/", Method: "GET", Handler: "RootHandler"},
	{Path: "/swagger/", Method: "GET", Handler: "SwaggerUI"},
	{Path: "/swagger/openapi.json", Method: "GET", Handler: "SwaggerSpec"},
	{Path: "/health", Method: "GET", Handler: "HealthHandler"},
	{Path: "/metrics", Method: "GET", Handler: "PrometheusMetrics"},
	{Path: "/api/auth/login", Method: "POST", Handler: "Login"},
	{Path: "/api/auth/register", Method: "POST", Handler: "Register"},
	{Path: "/api/auth/forgot-password", Method: "POST", Handler: "ForgotPassword"},
	{Path: "/api/auth/reset-password", Method: "POST", Handler: "ResetPassword"},
	{Path: "/api/auth/change-password", Method: "POST", Handler: "ChangePassword"},
	{Path: "/api/auth/logout", Method: "POST", Handler: "Logout"},
	{Path: "/api/auth/refresh-token", Method: "POST", Handler: "RefreshToken"},
	{Path: "/api/profile", Method: "GET", Handler: "GetProfile"},
	{Path: "/api/profile", Method: "PATCH", Handler: "UpdateProfile"},
	{Path: "/api/profile/password", Method: "PATCH", Handler: "ProfileChangePassword"},
	{Path: "/api/profile", Method: "DELETE", Handler: "DeactivateAccount"},
}

// GetRoutes returns all available routes
func GetRoutes() []Route {
	return availableRoutes // ← Update reference
}

// PrintRoutes prints all available routes
func PrintRoutes() {
	fmt.Println("\n📍 Available Routes:")
	fmt.Println("─────────────────────────────────────────")
	for i, r := range availableRoutes { // ← Update reference
		fmt.Printf("%d. [%s] %s | Handler: %s\n", i+1, r.Method, r.Path, r.Handler)
	}
	fmt.Printf("─────────────────────────────────────────\n")
	fmt.Printf("Total Routes: %d\n\n", len(availableRoutes)) // ← Update reference
}

// RegisterRoutes mendaftarkan semua HTTP routes
func RegisterRoutes(mux *http.ServeMux, db *database.DB) {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	tokenBlacklistRepo := repository.NewTokenBlacklistRepository(db)
	rateLimitRepo := repository.NewRateLimitRepository(db)
	passwordResetTokenRepo := repository.NewPasswordResetTokenRepository(db)
	passwordResetMailer := mailer.NewSMTPMailerFromEnv()
	tokenManager := security.NewJWTManager()

	// Initialize services
	authService := service.NewAuthService(userRepo, tokenBlacklistRepo, passwordResetTokenRepo, passwordResetMailer, tokenManager)
	profileService := service.NewProfileService(userRepo)
	rateLimiter := service.NewRateLimiterService(rateLimitRepo, time.Minute)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(rateLimiter)

	// Register route groups (sekarang OK!)
	routes.RegisterAuthRoutes(mux, authService, tokenManager, rateLimitMiddleware) // ← No conflict now!
	routes.RegisterProfileRoutes(mux, profileService, tokenManager, rateLimitMiddleware)

	// Public routes (health, root)
	mux.HandleFunc("GET /swagger/", SwaggerUIHandler())
	mux.HandleFunc("GET /swagger/openapi.json", SwaggerSpecHandler())
	mux.HandleFunc("GET /health", HealthHandler(db))
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("GET /", RootHandler())
}
