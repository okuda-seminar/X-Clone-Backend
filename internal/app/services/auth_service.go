package services

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	secretKey []byte
	logger    *slog.Logger
}

func NewAuthService(secretKey string) *AuthService {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return &AuthService{secretKey: []byte(secretKey), logger: logger}
}

// GenerateJWT generates a JWT with user ID and username
func (s *AuthService) GenerateJWT(ID string, username string) (string, error) {
	// Set payload (claims)
	claims := jwt.MapClaims{
		"sub":       ID,
		"username":  username,
		"exp":       time.Now().Add(time.Hour * 1).Unix(),       // Token expires after 1 hour
		"token_exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // Long-term expiration for 30 days (used for server-side validation)
	}

	// Set header & payload
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the JWT
	signedToken, err := unsignedToken.SignedString(s.secretKey)
	if err != nil {
		s.logger.Error("Failed to sign JWT", "error", err)
		return "", err
	}

	// Log the generated JWT
	s.logger.Info("Successfully generated JWT", "JWT", signedToken)
	return signedToken, nil
}

func (s *AuthService) ValidateJWT(token string) (jwt.MapClaims, error) {
	// Parse the JWT token and verify the signature
	parsedToken, err := jwt.Parse(token, func(unsignedToken *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	// If an error occurred during token parsing
	if err != nil {
		// Check if the error is due to expiration
		validationError, isValidationErr := err.(*jwt.ValidationError)
		if isValidationErr && (validationError.Errors&jwt.ValidationErrorExpired != 0) {
			s.logger.Warn("Token has expired", "error", err)
			return nil, fmt.Errorf("token has expired")
		}
		s.logger.Error("Invalid token", "error", err)
		return nil, fmt.Errorf("invalid token")
	}

	// Extract and cast the claims from the token
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Error("Failed to parse claims from token", "token", token)
		return nil, fmt.Errorf("failed to parse claims")
	}

	// Log and return the claims if the token is valid
	s.logger.Info("Token validated successfully", "claims", claims)
	return claims, nil
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
