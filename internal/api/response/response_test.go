package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON_StatusAndContentType(t *testing.T) {
	tests := []struct {
		name   string
		status int
		data   any
	}{
		{
			name:   "200 with map data",
			status: http.StatusOK,
			data:   map[string]string{"name": "test"},
		},
		{
			name:   "201 with struct data",
			status: http.StatusCreated,
			data:   struct{ ID int }{ID: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			JSON(w, tt.status, tt.data)

			if w.Code != tt.status {
				t.Errorf("status = %d, want %d", w.Code, tt.status)
			}

			ct := w.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("Content-Type = %q, want %q", ct, "application/json")
			}

			var env Envelope
			if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}
			if env.Data == nil {
				t.Error("envelope.Data is nil, expected data")
			}
			if env.Error != nil {
				t.Error("envelope.Error is non-nil, expected nil")
			}
		})
	}
}

func TestErr_Structure(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		code       string
		message    string
		callHelper func(w http.ResponseWriter)
	}{
		{
			name:    "bad request via Err",
			status:  http.StatusBadRequest,
			code:    "bad_request",
			message: "invalid input",
			callHelper: func(w http.ResponseWriter) {
				Err(w, http.StatusBadRequest, "bad_request", "invalid input")
			},
		},
		{
			name:    "bad request via helper",
			status:  http.StatusBadRequest,
			code:    "bad_request",
			message: "missing field",
			callHelper: func(w http.ResponseWriter) {
				BadRequest(w, "missing field")
			},
		},
		{
			name:    "unauthorized",
			status:  http.StatusUnauthorized,
			code:    "unauthorized",
			message: "token expired",
			callHelper: func(w http.ResponseWriter) {
				Unauthorized(w, "token expired")
			},
		},
		{
			name:    "forbidden",
			status:  http.StatusForbidden,
			code:    "forbidden",
			message: "access denied",
			callHelper: func(w http.ResponseWriter) {
				Forbidden(w, "access denied")
			},
		},
		{
			name:    "not found",
			status:  http.StatusNotFound,
			code:    "not_found",
			message: "resource not found",
			callHelper: func(w http.ResponseWriter) {
				NotFound(w, "resource not found")
			},
		},
		{
			name:    "internal error",
			status:  http.StatusInternalServerError,
			code:    "internal_error",
			message: "An unexpected error occurred",
			callHelper: func(w http.ResponseWriter) {
				InternalError(w)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.callHelper(w)

			if w.Code != tt.status {
				t.Errorf("status = %d, want %d", w.Code, tt.status)
			}

			ct := w.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("Content-Type = %q, want %q", ct, "application/json")
			}

			var env Envelope
			if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}
			if env.Data != nil {
				t.Error("envelope.Data should be nil for error responses")
			}
			if env.Error == nil {
				t.Fatal("envelope.Error is nil, expected error")
			}
			if env.Error.Code != tt.code {
				t.Errorf("error code = %q, want %q", env.Error.Code, tt.code)
			}
			if env.Error.Message != tt.message {
				t.Errorf("error message = %q, want %q", env.Error.Message, tt.message)
			}
		})
	}
}

func TestJSONWithMeta(t *testing.T) {
	w := httptest.NewRecorder()
	meta := &Meta{Page: 1, PerPage: 20, Total: 100, HasMore: true}
	JSONWithMeta(w, http.StatusOK, []string{"a", "b"}, meta)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var env Envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if env.Meta == nil {
		t.Fatal("envelope.Meta is nil, expected meta")
	}
	if env.Meta.Page != 1 {
		t.Errorf("meta.Page = %d, want 1", env.Meta.Page)
	}
	if env.Meta.Total != 100 {
		t.Errorf("meta.Total = %d, want 100", env.Meta.Total)
	}
}
