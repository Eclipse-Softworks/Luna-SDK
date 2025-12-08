// Package errors provides error types for the Luna SDK.
package errors

import "fmt"

// ErrorCode constants
const (
	CodeAuthInvalidKey             = "LUNA_ERR_AUTH_INVALID_KEY"
	CodeAuthTokenExpired           = "LUNA_ERR_AUTH_TOKEN_EXPIRED"
	CodeAuthInsufficientPermission = "LUNA_ERR_AUTH_INSUFFICIENT_PERMISSIONS"
	CodeRateLimitExceeded          = "LUNA_ERR_RATE_LIMIT_EXCEEDED"
	CodeResourceNotFound           = "LUNA_ERR_RESOURCE_NOT_FOUND"
	CodeResourceConflict           = "LUNA_ERR_RESOURCE_CONFLICT"
	CodeValidationFailed           = "LUNA_ERR_VALIDATION_FAILED"
	CodeValidationInvalidParam     = "LUNA_ERR_VALIDATION_INVALID_PARAMETER"
	CodeNetworkTimeout             = "LUNA_ERR_NETWORK_TIMEOUT"
	CodeNetworkConnection          = "LUNA_ERR_NETWORK_CONNECTION"
	CodeServerInternal             = "LUNA_ERR_SERVER_INTERNAL"
	CodeServerUnavailable          = "LUNA_ERR_SERVER_UNAVAILABLE"
)

// retryableCodes contains error codes that are safe to retry
var retryableCodes = map[string]bool{
	CodeRateLimitExceeded: true,
	CodeNetworkTimeout:    true,
	CodeNetworkConnection: true,
	CodeServerInternal:    true,
	CodeServerUnavailable: true,
}

// Error is the base error type for Luna SDK
type Error struct {
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Status    int                    `json:"status"`
	RequestID string                 `json:"request_id"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

func (e *Error) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("[%s] %s (Request ID: %s)", e.Code, e.Message, e.RequestID)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// DocsURL returns the documentation URL for this error
func (e *Error) DocsURL() string {
	return fmt.Sprintf("https://docs.eclipse.dev/luna/errors#%s", e.Code)
}

// Retryable returns true if this error is safe to retry
func (e *Error) Retryable() bool {
	return retryableCodes[e.Code]
}

// AuthenticationError indicates authentication failure
type AuthenticationError struct {
	BaseError *Error
}

func (e *AuthenticationError) Error() string {
	return e.BaseError.Error()
}

// AuthorizationError indicates authorization failure
type AuthorizationError struct {
	BaseError *Error
}

func (e *AuthorizationError) Error() string {
	return e.BaseError.Error()
}

// ValidationError indicates validation failure
type ValidationError struct {
	BaseError *Error
}

func (e *ValidationError) Error() string {
	return e.BaseError.Error()
}

// RateLimitError indicates rate limit exceeded
type RateLimitError struct {
	BaseError  *Error
	RetryAfter int
}

func (e *RateLimitError) Error() string {
	return e.BaseError.Error()
}

// NetworkError indicates network-related errors
type NetworkError struct {
	BaseError *Error
}

func (e *NetworkError) Error() string {
	return e.BaseError.Error()
}

// NotFoundError indicates resource not found
type NotFoundError struct {
	BaseError *Error
}

func (e *NotFoundError) Error() string {
	return e.BaseError.Error()
}

// ConflictError indicates resource conflict
type ConflictError struct {
	BaseError *Error
}

func (e *ConflictError) Error() string {
	return e.BaseError.Error()
}

// ServerError indicates server-side errors
type ServerError struct {
	BaseError *Error
}

func (e *ServerError) Error() string {
	return e.BaseError.Error()
}

// New creates a new Error with the given parameters
func New(code, message string, status int, requestID string, details map[string]interface{}) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Status:    status,
		RequestID: requestID,
		Details:   details,
	}
}

// FromResponse creates an appropriate error type from an API response
func FromResponse(status int, code, message, requestID string, details map[string]interface{}, retryAfter int) error {
	base := &Error{
		Code:      code,
		Message:   message,
		Status:    status,
		RequestID: requestID,
		Details:   details,
	}

	switch status {
	case 400:
		return &ValidationError{BaseError: base}
	case 401:
		return &AuthenticationError{BaseError: base}
	case 403:
		return &AuthorizationError{BaseError: base}
	case 404:
		return &NotFoundError{BaseError: base}
	case 409:
		return &ConflictError{BaseError: base}
	case 429:
		return &RateLimitError{BaseError: base, RetryAfter: retryAfter}
	default:
		if status >= 500 {
			return &ServerError{BaseError: base}
		}
		return base
	}
}
