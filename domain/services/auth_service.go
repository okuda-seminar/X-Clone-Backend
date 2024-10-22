package services

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	secretKey []byte
}

func NewAuthService(secretKey string) *AuthService {
	return &AuthService{secretKey: []byte(secretKey)}
}

func (s *AuthService) GenerateJWT(ID int, Username string) (string, error) {
	// Set payload (claims)
	claims := jwt.MapClaims{
		"sub":       ID,
		"username":  Username,
		"exp":       time.Now().Add(time.Hour * 1).Unix(),       // Token expires after 1 hour
		"token_exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // Long-term expiration for 30 days (used for server-side validation)
	}

	//set header & payload
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the JWT
	signedToken, err := unsignedToken.SignedString(s.secretKey)
	if err != nil {
		fmt.Println("JWT signing error:", err)
		return "", err
	}

	// Display the generated JWT
	fmt.Println("JWT:", string(signedToken))
	return signedToken, nil
}
