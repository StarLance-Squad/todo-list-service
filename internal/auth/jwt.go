package auth

import (
	"errors"
	"github.com/golang-jwt/jwt"
)

// JwtCustomClaims holds the claims for the JWT token.
type JwtCustomClaims struct {
	UserID uint `json:"user_id"`
	Admin  bool `json:"admin"`
	jwt.StandardClaims
}

// GenerateToken creates a JWT token with the specified username, userID, and secret.
// It returns an error if the username or secret is empty.
func GenerateToken(username string, userID uint, secret string) (string, error) {
	if username == "" {
		return "", errors.New("username cannot be empty")
	}
	if secret == "" {
		return "", errors.New("secret cannot be empty")
	}

	claims := &JwtCustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			Subject: username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
