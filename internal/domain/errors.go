package domain

import (
	"errors"
	"fmt"
)

// Domain errors
var (
	ErrUserNotFound       = fmt.Errorf("user not found")
	ErrUserAlreadyExists  = fmt.Errorf("user already exists")
	ErrInvalidCredentials = fmt.Errorf("invalid email/username or password")
	ErrInvalidEmail       = fmt.Errorf("invalid email format")
	ErrInvalidResetToken  = fmt.Errorf("invalid or expired reset token")
	ErrPasswordTooShort   = fmt.Errorf("password must be at least 6 characters")
	ErrInvalidRequest     = fmt.Errorf("invalid request")
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
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
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

type InternalError struct {
	Err error
}

func (e InternalError) Error() string {
	return ErrInternal.Error()
}

func (e InternalError) Unwrap() error {
	return e.Err
}

func NewInternalError(err error) error {
	if err == nil {
		return ErrInternal
	}
	return InternalError{Err: err}
}

func DebugError(err error) string {
	var internalErr InternalError
	if errors.As(err, &internalErr) && internalErr.Err != nil {
		return internalErr.Err.Error()
	}

	if err != nil {
		return err.Error()
	}
	return ""
}

// Mapping error ke HTTP status code
func ErrorToStatusCode(err error) int {
	if err == nil {
		return 200
	}

	var businessErr BusinessError
	if errors.As(err, &businessErr) {
		err = businessErr.Err
	}

	var internalErr InternalError
	if errors.As(err, &internalErr) {
		return 500
	}

	switch err {
	case ErrUserNotFound, ErrBadgeNotFound, ErrDictionaryNotFound:
		return 404
	case ErrInvalidRequest, ErrInvalidEmail, ErrInvalidResetToken, ErrPasswordTooShort:
		return 400
	case ErrInvalidCredentials:
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
