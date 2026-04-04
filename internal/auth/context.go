package auth

import (
	"context"

	"github.com/google/uuid"
)

// AuthType represents the method used to authenticate a request.
type AuthType string

const (
	// AuthTypeJWT indicates authentication via a JWT token.
	AuthTypeJWT AuthType = "jwt"
	// AuthTypeAPIKey indicates authentication via an API key.
	AuthTypeAPIKey AuthType = "apikey"
)

type ctxKey string

const (
	ctxUserID   ctxKey = "auth_user_id"
	ctxOrgID    ctxKey = "auth_org_id"
	ctxRole     ctxKey = "auth_role"
	ctxAuthType ctxKey = "auth_type"
	ctxAPIKeyID ctxKey = "auth_api_key_id"
)

// WithAuthContext stores authentication information in the context.
func WithAuthContext(ctx context.Context, userID, orgID uuid.UUID, role string, authType AuthType, apiKeyID *uuid.UUID) context.Context {
	ctx = context.WithValue(ctx, ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxOrgID, orgID)
	ctx = context.WithValue(ctx, ctxRole, role)
	ctx = context.WithValue(ctx, ctxAuthType, authType)
	ctx = context.WithValue(ctx, ctxAPIKeyID, apiKeyID)
	return ctx
}

// UserIDFromContext retrieves the authenticated user's ID from the context.
func UserIDFromContext(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(ctxUserID).(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}

// OrgIDFromContext retrieves the authenticated user's organization ID from the context.
func OrgIDFromContext(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(ctxOrgID).(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}

// RoleFromContext retrieves the authenticated user's role from the context.
func RoleFromContext(ctx context.Context) string {
	if role, ok := ctx.Value(ctxRole).(string); ok {
		return role
	}
	return ""
}

// AuthTypeFromContext retrieves the authentication type from the context.
func AuthTypeFromContext(ctx context.Context) AuthType {
	if at, ok := ctx.Value(ctxAuthType).(AuthType); ok {
		return at
	}
	return ""
}

// APIKeyIDFromContext retrieves the API key ID from the context, if present.
func APIKeyIDFromContext(ctx context.Context) *uuid.UUID {
	if id, ok := ctx.Value(ctxAPIKeyID).(*uuid.UUID); ok {
		return id
	}
	return nil
}
