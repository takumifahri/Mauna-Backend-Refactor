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
	{Path: "/api/auth/verify-registration", Method: "POST", Handler: "VerifyRegistration"},
	{Path: "/api/auth/forgot-password", Method: "POST", Handler: "ForgotPassword"},
	{Path: "/api/auth/reset-password", Method: "POST", Handler: "ResetPassword"},
	{Path: "/api/auth/change-password", Method: "POST", Handler: "ChangePassword"},
	{Path: "/api/auth/logout", Method: "POST", Handler: "Logout"},
	{Path: "/api/auth/refresh-token", Method: "POST", Handler: "RefreshToken"},
	{Path: "/api/profile", Method: "GET", Handler: "GetProfile"},
	{Path: "/api/profile", Method: "PATCH", Handler: "UpdateProfile"},
	{Path: "/api/profile/avatar", Method: "POST", Handler: "UploadAvatar"},
	{Path: "/api/profile/password", Method: "PATCH", Handler: "ProfileChangePassword"},
	{Path: "/api/profile", Method: "DELETE", Handler: "DeactivateAccount"},
	{Path: "/api/admin/users", Method: "GET", Handler: "AdminListUsers"},
	{Path: "/api/admin/users", Method: "POST", Handler: "AdminCreateUser"},
	{Path: "/api/admin/users/{id}", Method: "GET", Handler: "AdminGetUser"},
	{Path: "/api/admin/users/{id}", Method: "PATCH", Handler: "AdminUpdateUser"},
	{Path: "/api/admin/users/{id}/soft-delete", Method: "PATCH", Handler: "AdminSoftDeleteUser"},
	{Path: "/api/admin/users/{id}", Method: "DELETE", Handler: "AdminDeleteUser"},
	{Path: "/api/admin/users/{id}/restore", Method: "PATCH", Handler: "AdminRestoreUser"},
	{Path: "/api/admin/badges", Method: "GET", Handler: "AdminListBadges"},
	{Path: "/api/admin/badges", Method: "POST", Handler: "AdminCreateBadge"},
	{Path: "/api/admin/badges/{id}", Method: "GET", Handler: "AdminGetBadge"},
	{Path: "/api/admin/badges/{id}", Method: "PATCH", Handler: "AdminUpdateBadge"},
	{Path: "/api/admin/badges/{id}/soft-delete", Method: "PATCH", Handler: "AdminSoftDeleteBadge"},
	{Path: "/api/admin/badges/{id}", Method: "DELETE", Handler: "AdminDeleteBadge"},
	{Path: "/api/admin/badges/{id}/restore", Method: "PATCH", Handler: "AdminRestoreBadge"},
	{Path: "/api/admin/dictionary", Method: "GET", Handler: "AdminListDictionary"},
	{Path: "/api/admin/dictionary", Method: "POST", Handler: "AdminCreateDictionary"},
	{Path: "/api/admin/dictionary/{id}", Method: "GET", Handler: "AdminGetDictionary"},
	{Path: "/api/admin/dictionary/{id}", Method: "PATCH", Handler: "AdminUpdateDictionary"},
	{Path: "/api/admin/dictionary/{id}/soft-delete", Method: "PATCH", Handler: "AdminSoftDeleteDictionary"},
	{Path: "/api/admin/dictionary/{id}", Method: "DELETE", Handler: "AdminDeleteDictionary"},
	{Path: "/api/admin/dictionary/{id}/restore", Method: "PATCH", Handler: "AdminRestoreDictionary"},
	{Path: "/api/admin/levels", Method: "GET", Handler: "AdminListLevels"},
	{Path: "/api/admin/levels", Method: "POST", Handler: "AdminCreateLevel"},
	{Path: "/api/admin/levels/{id}", Method: "GET", Handler: "AdminGetLevel"},
	{Path: "/api/admin/levels/{id}", Method: "PATCH", Handler: "AdminUpdateLevel"},
	{Path: "/api/admin/levels/{id}/soft-delete", Method: "PATCH", Handler: "AdminSoftDeleteLevel"},
	{Path: "/api/admin/levels/{id}", Method: "DELETE", Handler: "AdminDeleteLevel"},
	{Path: "/api/admin/levels/{id}/restore", Method: "PATCH", Handler: "AdminRestoreLevel"},
	{Path: "/api/admin/sublevels", Method: "GET", Handler: "AdminListSubLevels"},
	{Path: "/api/admin/sublevels", Method: "POST", Handler: "AdminCreateSubLevel"},
	{Path: "/api/admin/sublevels/{id}", Method: "GET", Handler: "AdminGetSubLevel"},
	{Path: "/api/admin/sublevels/{id}", Method: "PATCH", Handler: "AdminUpdateSubLevel"},
	{Path: "/api/admin/sublevels/{id}/soft-delete", Method: "PATCH", Handler: "AdminSoftDeleteSubLevel"},
	{Path: "/api/admin/sublevels/{id}", Method: "DELETE", Handler: "AdminDeleteSubLevel"},
	{Path: "/api/admin/sublevels/{id}/restore", Method: "PATCH", Handler: "AdminRestoreSubLevel"},
	{Path: "/api/admin/soal", Method: "GET", Handler: "AdminListSoal"},
	{Path: "/api/admin/soal", Method: "POST", Handler: "AdminCreateSoal"},
	{Path: "/api/admin/soal/{id}", Method: "GET", Handler: "AdminGetSoal"},
	{Path: "/api/admin/soal/{id}", Method: "PATCH", Handler: "AdminUpdateSoal"},
	{Path: "/api/admin/soal/{id}/soft-delete", Method: "PATCH", Handler: "AdminSoftDeleteSoal"},
	{Path: "/api/admin/soal/{id}", Method: "DELETE", Handler: "AdminDeleteSoal"},
	{Path: "/api/admin/soal/{id}/restore", Method: "PATCH", Handler: "AdminRestoreSoal"},
	{Path: "/api/admin/progress", Method: "GET", Handler: "AdminListProgress"},
	{Path: "/api/admin/progress/{id}", Method: "GET", Handler: "AdminGetProgress"},
	{Path: "/api/admin/progress/{id}/soft-delete", Method: "PATCH", Handler: "AdminSoftDeleteProgress"},
	{Path: "/api/admin/progress/{id}", Method: "DELETE", Handler: "AdminDeleteProgress"},
	{Path: "/api/admin/progress/{id}/restore", Method: "PATCH", Handler: "AdminRestoreProgress"},
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
	pendingRegistrationRepo := repository.NewPendingRegistrationRepository(db)
	managementUsersRepo := repository.NewManagementUsersRepository(db)
	managementBadgesRepo := repository.NewManagementBadgesRepository(db)
	managementDictionaryRepo := repository.NewManagementDictionaryRepository(db)
	managementLevelsRepo := repository.NewManagementLevelsRepository(db)
	managementSubLevelsRepo := repository.NewManagementSubLevelsRepository(db)
	managementSoalRepo := repository.NewManagementSoalRepository(db)
	managementProgressRepo := repository.NewManagementProgressRepository(db)
	passwordResetMailer := mailer.NewSMTPMailerFromEnv()
	tokenManager := security.NewJWTManager()

	// Initialize services
	authService := service.NewAuthService(userRepo, tokenBlacklistRepo, passwordResetTokenRepo, pendingRegistrationRepo, passwordResetMailer, tokenManager)
	profileService := service.NewProfileService(userRepo)
	adminUsersService := service.NewManagementUsersService(managementUsersRepo)
	adminBadgesService := service.NewManagementBadgesService(managementBadgesRepo)
	adminDictionaryService := service.NewManagementDictionaryService(managementDictionaryRepo)
	adminLevelsService := service.NewManagementLevelsService(managementLevelsRepo)
	adminSubLevelsService := service.NewManagementSubLevelsService(managementSubLevelsRepo)
	adminSoalService := service.NewManagementSoalService(managementSoalRepo)
	adminProgressService := service.NewManagementProgressService(managementProgressRepo)
	rateLimiter := service.NewRateLimiterService(rateLimitRepo, time.Minute)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(rateLimiter)

	// Register route groups (sekarang OK!)
	routes.RegisterAuthRoutes(mux, authService, tokenManager, rateLimitMiddleware) // ← No conflict now!
	routes.RegisterProfileRoutes(mux, profileService, tokenManager, rateLimitMiddleware)
	routes.RegisterAdminUsersRoutes(mux, adminUsersService, tokenManager, rateLimitMiddleware)
	routes.RegisterAdminBadgesRoutes(mux, adminBadgesService, tokenManager, rateLimitMiddleware)
	routes.RegisterAdminDictionaryRoutes(mux, adminDictionaryService, tokenManager, rateLimitMiddleware)
	routes.RegisterAdminLevelsRoutes(mux, adminLevelsService, adminSubLevelsService, tokenManager, rateLimitMiddleware)
	routes.RegisterAdminSoalRoutes(mux, adminSoalService, tokenManager, rateLimitMiddleware)
	routes.RegisterAdminProgressRoutes(mux, adminProgressService, tokenManager, rateLimitMiddleware)

	// Public routes (health, root)
	mux.HandleFunc("GET /swagger/", SwaggerUIHandler())
	mux.HandleFunc("GET /swagger/openapi.json", SwaggerSpecHandler())
	mux.HandleFunc("GET /health", HealthHandler(db))
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	mux.HandleFunc("GET /", RootHandler())
}
