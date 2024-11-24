package services

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

// TestGenerateJWT tests the JWT generation functionality of AuthService.
// It ensures that a token is generated without errors and is not empty.
func TestGenerateJWT(t *testing.T) {
	secretKey := "test_secret_key"
	authService := NewAuthService(secretKey)

	userID := "1"
	username := "test_user"

	signedToken, err := authService.GenerateJWT(userID, username)

	assert.NoError(t, err)
	assert.NotEmpty(t, signedToken)
}

// TestValidateJWT tests the validation of a valid JWT in AuthService.
// It checks if the generated token's claims match the expected values.
func TestValidateJWT(t *testing.T) {
	secretKey := "test_secret_key"
	authService := NewAuthService(secretKey)

	userID := "1"
	username := "test_user"

	signedToken, err := authService.GenerateJWT(userID, username)
	assert.NoError(t, err)

	claims, err := authService.ValidateJWT(signedToken)
	assert.NoError(t, err)

	assert.Equal(t, userID, claims["sub"])
	assert.Equal(t, username, claims["username"])
}

// TestExpiredJWT tests the validation of an expired JWT.
// It verifies that the error returned is "token has expired".
func TestExpiredJWT(t *testing.T) {
	secretKey := "test_secret_key"
	authService := NewAuthService(secretKey)

	expiredToken := generateExpiredJWT(secretKey)

	_, err := authService.ValidateJWT(expiredToken)
	assert.EqualError(t, err, "token has expired")
}

// TestInvalidSignatureJWT tests the validation of a JWT with an invalid signature.
// It ensures that the error returned is "invalid token".
func TestInvalidSignatureJWT(t *testing.T) {
	secretKey := "test_secret_key"
	authService := NewAuthService(secretKey)

	userID := "1"
	username := "test_user"

	signedToken, err := authService.GenerateJWT(userID, username)
	assert.NoError(t, err)

	invalidToken := signedToken[:len(signedToken)-1] + "x"

	_, err = authService.ValidateJWT(invalidToken)
	assert.EqualError(t, err, "invalid token")
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
