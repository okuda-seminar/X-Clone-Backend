package services

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// TestGenerateJWT tests the JWT generation functionality of AuthService.
func TestGenerateJWT(t *testing.T) {
	secretKey := "test_secret_key"
	authService := NewAuthService(secretKey)

	userID := uuid.New()
	username := "test_user"

	signedToken, err := authService.GenerateJWT(userID, username)

	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	if signedToken == "" {
		t.Fatalf("Expected a valid token, but got an empty string")
	}
}

// TestValidateJWT tests the validation of a valid JWT in AuthService.
func TestValidateJWT(t *testing.T) {
	secretKey := "test_secret_key"
	authService := NewAuthService(secretKey)

	userID := uuid.New()
	username := "test_user"

	signedToken, err := authService.GenerateJWT(userID, username)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	claims, err := authService.ValidateJWT(signedToken)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if claims["sub"] != userID.String() {
		t.Errorf("Expected user ID %v, but got %v", userID, claims["sub"])
	}
	if claims["username"] != username {
		t.Errorf("Expected username %v, but got %v", username, claims["username"])
	}
}

// TestExpiredJWT tests the validation of an expired JWT.
func TestExpiredJWT(t *testing.T) {
	secretKey := "test_secret_key"
	authService := NewAuthService(secretKey)

	expiredToken := generateExpiredJWT(secretKey)

	_, err := authService.ValidateJWT(expiredToken)
	if err == nil || err.Error() != "token has expired" {
		t.Errorf("Expected error 'token has expired', but got: %v", err)
	}
}

// TestInvalidSignatureJWT tests the validation of a JWT with an invalid signature.
func TestInvalidSignatureJWT(t *testing.T) {
	secretKey := "test_secret_key"
	authService := NewAuthService(secretKey)

	userID := uuid.New()
	username := "test_user"

	signedToken, err := authService.GenerateJWT(userID, username)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	invalidToken := signedToken[:len(signedToken)-1] + "x"

	_, err = authService.ValidateJWT(invalidToken)
	if err == nil || err.Error() != "invalid token" {
		t.Errorf("Expected error 'invalid token', but got: %v", err)
	}
}

// generateExpiredJWT generates an expired JWT for testing purposes.
func generateExpiredJWT(secretKey string) string {
	claims := jwt.MapClaims{
		"sub":       1,
		"username":  "test_user",
		"exp":       time.Now().Add(-time.Hour).Unix(),
		"token_exp": time.Now().Add(-time.Hour).Unix(),
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := unsignedToken.SignedString([]byte(secretKey))
	return signedToken
}
