package dto

// Response wrapper untuk semua endpoint
type Response struct {
    Status    string      `json:"status"`
    Message   string      `json:"message,omitempty"`
    Data      interface{} `json:"data,omitempty"`
    Timestamp string      `json:"timestamp"`
}

// Error response
type ErrorResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    Error   string `json:"error"`
}