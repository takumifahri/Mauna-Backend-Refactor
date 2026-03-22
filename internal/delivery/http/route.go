package http

import (
    "fmt"
    "net/http"

    "REFACTORING_MAUNA/pkg/database"
)

// Route represents a route definition
type Route struct {
    Path    string
    Method  string
    Handler string
}

var routes = []Route{
    {Path: "/", Method: "GET", Handler: "RootHandler"},
    {Path: "/health", Method: "GET", Handler: "HealthHandler"},
}

// GetRoutes returns all available routes
func GetRoutes() []Route {
    return routes
}

// PrintRoutes prints all available routes
func PrintRoutes() {
    fmt.Println("\n📍 Available Routes:")
    fmt.Println("─────────────────────────────────────────")
    for i, r := range routes {
        fmt.Printf("%d. [%s] %s | Handler: %s\n", i+1, r.Method, r.Path, r.Handler)
    }
    fmt.Printf("─────────────────────────────────────────\n")
    fmt.Printf("Total Routes: %d\n\n", len(routes))
}

// RegisterRoutes mendaftarkan semua HTTP routes
func RegisterRoutes(mux *http.ServeMux, db *database.DB) {
    // Health check
    mux.HandleFunc("/health", HealthHandler(db))

    // Root/Info endpoint
    mux.HandleFunc("/", RootHandler())
}