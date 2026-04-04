package response

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Data  any    `json:"data,omitempty"`
	Error *Error `json:"error,omitempty"`
	Meta  *Meta  `json:"meta,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page       int  `json:"page,omitempty"`
	PerPage    int  `json:"per_page,omitempty"`
	Total      int  `json:"total,omitempty"`
	HasMore    bool `json:"has_more,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Envelope{Data: data})
}

func JSONWithMeta(w http.ResponseWriter, status int, data any, meta *Meta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Envelope{Data: data, Meta: meta})
}

func Err(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Envelope{
		Error: &Error{
			Code:    code,
			Message: message,
		},
	})
}

func BadRequest(w http.ResponseWriter, message string) {
	Err(w, http.StatusBadRequest, "bad_request", message)
}

func Unauthorized(w http.ResponseWriter, message string) {
	Err(w, http.StatusUnauthorized, "unauthorized", message)
}

func Forbidden(w http.ResponseWriter, message string) {
	Err(w, http.StatusForbidden, "forbidden", message)
}

func NotFound(w http.ResponseWriter, message string) {
	Err(w, http.StatusNotFound, "not_found", message)
}

func Conflict(w http.ResponseWriter, message string) {
	Err(w, http.StatusConflict, "conflict", message)
}

func TooManyRequests(w http.ResponseWriter, message string) {
	Err(w, http.StatusTooManyRequests, "rate_limit_exceeded", message)
}

func InternalError(w http.ResponseWriter) {
	Err(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred")
}
