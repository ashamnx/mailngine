package middleware

import (
	"net/http"

	"github.com/mailngine/mailngine/internal/api/response"
	"github.com/mailngine/mailngine/internal/auth"
)

// roleHierarchy defines the privilege levels for each role.
// Higher values indicate more privileges.
var roleHierarchy = map[string]int{
	"viewer": 0,
	"member": 1,
	"admin":  2,
	"owner":  3,
}

// RequireRole returns middleware that enforces a minimum role level.
// The authenticated user's role (from context) must meet or exceed the specified minimum.
// For API key authentication, the permission is mapped to a role via permissionToRole
// which is already set during authentication.
func RequireRole(minRole string) func(http.Handler) http.Handler {
	minLevel, ok := roleHierarchy[minRole]
	if !ok {
		minLevel = 0
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := auth.RoleFromContext(r.Context())
			userLevel, ok := roleHierarchy[role]
			if !ok {
				response.Forbidden(w, "insufficient permissions")
				return
			}

			if userLevel < minLevel {
				response.Forbidden(w, "insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
