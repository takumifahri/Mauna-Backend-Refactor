package domain

import "fmt"

// Domain errors
var (
    ErrUserNotFound       = fmt.Errorf("user not found")
    ErrUserAlreadyExists  = fmt.Errorf("user already exists")
    ErrInvalidCredentials = fmt.Errorf("invalid email/username or password")
    ErrInvalidEmail       = fmt.Errorf("invalid email format")
    ErrPasswordTooShort   = fmt.Errorf("password must be at least 6 characters")
    ErrUnauthorized       = fmt.Errorf("unauthorized")
    ErrForbidden          = fmt.Errorf("forbidden")
    ErrInternal           = fmt.Errorf("internal server error")
    ErrBadgeNotFound      = fmt.Errorf("badge not found")
    ErrDictionaryNotFound = fmt.Errorf("dictionary not found")
    ErrLevelNotFound      = fmt.Errorf("level not found")
    ErrQuestionNotFound   = fmt.Errorf("question not found")
    ErrProgressNotFound   = fmt.Errorf("progress not found")
)

// Custom error type untuk response
type ErrorDetail struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}

// BusinessError adalah error yang bisa di-handle
type BusinessError struct {
    Code    string
    Message string
    Err     error
}

func (e BusinessError) Error() string {
    return e.Message
}

// Helper functions
func NewBusinessError(code, message string, err error) BusinessError {
    return BusinessError{
        Code:    code,
        Message: message,
        Err:     err,
    }
}

// Mapping error ke HTTP status code
func ErrorToStatusCode(err error) int {
    if err == nil {
        return 200
    }

    switch err {
    case ErrUserNotFound, ErrBadgeNotFound, ErrDictionaryNotFound:
        return 404
    case ErrInvalidCredentials, ErrInvalidEmail:
        return 401
    case ErrUserAlreadyExists:
        return 409
    case ErrUnauthorized:
        return 401
    case ErrForbidden:
        return 403
    default:
        return 500
    }
}