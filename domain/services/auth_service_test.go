package services

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

// TestGenerateJWT tests the JWT generation functionality of AuthService.
// It verifies that a JWT is generated correctly with valid claims and
// is successfully parsed using the same secret key.
func TestGenerateJWT(t *testing.T) {
	secretKey := "test_secret_key"
	authService := NewAuthService(secretKey)

	userID := 1
	username := "test_user"

	// Generate the JWT for a specified user ID and username.
	signedToken, err := authService.GenerateJWT(userID, username)

	// Validate that no error occurred and the token is not empty.
	assert.NoError(t, err)
	assert.NotEmpty(t, signedToken)

	// Parse the generated JWT to verify its integrity.
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	// Ensure the token is valid and parsed correctly.
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	// Check the claims in the token.
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, float64(userID), claims["sub"])
	assert.Equal(t, username, claims["username"])
}
