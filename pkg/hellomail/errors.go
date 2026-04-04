package hellomail

import "fmt"

// APIError represents a structured error from the Hello Mail API.
type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Common error codes
var (
	ErrBadRequest       = &APIError{StatusCode: 400, Code: "bad_request", Message: "The request was invalid"}
	ErrUnauthorized     = &APIError{StatusCode: 401, Code: "unauthorized", Message: "Authentication required"}
	ErrForbidden        = &APIError{StatusCode: 403, Code: "forbidden", Message: "You don't have permission to perform this action"}
	ErrNotFound         = &APIError{StatusCode: 404, Code: "not_found", Message: "The requested resource was not found"}
	ErrConflict         = &APIError{StatusCode: 409, Code: "conflict", Message: "The resource already exists"}
	ErrRateLimited      = &APIError{StatusCode: 429, Code: "rate_limit_exceeded", Message: "Too many requests"}
	ErrInternalError    = &APIError{StatusCode: 500, Code: "internal_error", Message: "An unexpected error occurred"}
	ErrDomainNotVerified = &APIError{StatusCode: 422, Code: "domain_not_verified", Message: "The sending domain has not been verified"}
	ErrSuppressed       = &APIError{StatusCode: 422, Code: "recipient_suppressed", Message: "The recipient is on the suppression list"}
	ErrLimitExceeded    = &APIError{StatusCode: 429, Code: "limit_exceeded", Message: "Monthly email limit exceeded"}
)
