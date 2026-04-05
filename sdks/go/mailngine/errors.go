package mailngine

import (
	"errors"
	"fmt"
)

// APIError represents a structured error response from the Mailngine API.
type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("mailngine: %s: %s (status %d)", e.Code, e.Message, e.StatusCode)
}

// IsNotFound reports whether the error indicates a 404 Not Found response.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 404
	}
	return false
}

// IsRateLimited reports whether the error indicates a 429 Too Many Requests response.
func IsRateLimited(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 429
	}
	return false
}

// IsValidationError reports whether the error indicates a 400 Bad Request response.
func IsValidationError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 400
	}
	return false
}

// IsAuthenticationError reports whether the error indicates a 401 Unauthorized response.
func IsAuthenticationError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 401
	}
	return false
}
