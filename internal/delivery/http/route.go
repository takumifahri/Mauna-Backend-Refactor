package http

import (
    "fmt"
    "net/http"

    "REFACTORING_MAUNA/internal/delivery/http/routes"
    "REFACTORING_MAUNA/internal/repository"
    "REFACTORING_MAUNA/internal/service"
    "REFACTORING_MAUNA/pkg/database"
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
    {Path: "/health", Method: "GET", Handler: "HealthHandler"},
    {Path: "/api/auth/login", Method: "POST", Handler: "Login"},
    {Path: "/api/auth/register", Method: "POST", Handler: "Register"},
    {Path: "/api/auth/change-password", Method: "POST", Handler: "ChangePassword"},
    {Path: "/api/auth/logout", Method: "POST", Handler: "Logout"},
    {Path: "/api/auth/refresh-token", Method: "POST", Handler: "RefreshToken"},
}

// GetRoutes returns all available routes
func GetRoutes() []Route {
    return availableRoutes  // ← Update reference
}

// PrintRoutes prints all available routes
func PrintRoutes() {
    fmt.Println("\n📍 Available Routes:")
    fmt.Println("─────────────────────────────────────────")
    for i, r := range availableRoutes {  // ← Update reference
        fmt.Printf("%d. [%s] %s | Handler: %s\n", i+1, r.Method, r.Path, r.Handler)
    }
    fmt.Printf("─────────────────────────────────────────\n")
    fmt.Printf("Total Routes: %d\n\n", len(availableRoutes))  // ← Update reference
}

// RegisterRoutes mendaftarkan semua HTTP routes
func RegisterRoutes(mux *http.ServeMux, db *database.DB) {
    // Initialize repositories
    userRepo := repository.NewUserRepository(db)

    // Initialize services
    authService := service.NewAuthService(userRepo)

    // Register route groups (sekarang OK!)
    routes.RegisterAuthRoutes(mux, authService)  // ← No conflict now!

    // Public routes (health, root)
    mux.HandleFunc("GET /health", HealthHandler(db))
    mux.HandleFunc("GET /", RootHandler())
}