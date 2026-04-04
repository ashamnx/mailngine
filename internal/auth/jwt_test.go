package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestJWTManager_GenerateAndValidate(t *testing.T) {
	secret := "test-secret-key-for-jwt-testing"
	expiry := 15 * time.Minute
	mgr := NewJWTManager(secret, expiry)

	userID := uuid.New()
	orgID := uuid.New()
	role := "admin"

	token, err := mgr.Generate(userID, orgID, role)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if token == "" {
		t.Fatal("Generate() returned empty token")
	}

	claims, err := mgr.Validate(token)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.OrgID != orgID {
		t.Errorf("OrgID = %v, want %v", claims.OrgID, orgID)
	}
	if claims.Role != role {
		t.Errorf("Role = %q, want %q", claims.Role, role)
	}
	if claims.Issuer != "hellomail" {
		t.Errorf("Issuer = %q, want %q", claims.Issuer, "hellomail")
	}
	if claims.Subject != userID.String() {
		t.Errorf("Subject = %q, want %q", claims.Subject, userID.String())
	}
}

func TestJWTManager_ExpiredTokenRejected(t *testing.T) {
	secret := "test-secret-key-for-jwt-testing"
	// Use a negative expiry to produce an already-expired token.
	mgr := NewJWTManager(secret, -1*time.Hour)

	userID := uuid.New()
	orgID := uuid.New()

	token, err := mgr.Generate(userID, orgID, "member")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	_, err = mgr.Validate(token)
	if err == nil {
		t.Fatal("Validate() expected error for expired token, got nil")
	}
}

func TestJWTManager_InvalidSignatureRejected(t *testing.T) {
	mgr1 := NewJWTManager("secret-one", 15*time.Minute)
	mgr2 := NewJWTManager("secret-two", 15*time.Minute)

	userID := uuid.New()
	orgID := uuid.New()

	token, err := mgr1.Generate(userID, orgID, "member")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Validate with a different secret should fail.
	_, err = mgr2.Validate(token)
	if err == nil {
		t.Fatal("Validate() expected error for invalid signature, got nil")
	}
}

func TestJWTManager_MalformedTokenRejected(t *testing.T) {
	mgr := NewJWTManager("test-secret", 15*time.Minute)

	_, err := mgr.Validate("not-a-valid-jwt")
	if err == nil {
		t.Fatal("Validate() expected error for malformed token, got nil")
	}
}

func TestJWTManager_WrongSigningMethodRejected(t *testing.T) {
	mgr := NewJWTManager("test-secret", 15*time.Minute)

	// Create a token with an unexpected signing method (none).
	token := jwt.NewWithClaims(jwt.SigningMethodNone, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "hellomail",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
		UserID: uuid.New(),
		OrgID:  uuid.New(),
		Role:   "admin",
	})
	signed, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("failed to sign token with none method: %v", err)
	}

	_, err = mgr.Validate(signed)
	if err == nil {
		t.Fatal("Validate() expected error for wrong signing method, got nil")
	}
}
