package auth

import (
	"github.com/golang-jwt/jwt"
)

type JwtCustomClaims struct {
	UserID uint `json:"user_id"`
	Admin  bool `json:"admin"`
	jwt.StandardClaims
}

func GenerateToken(username string, userID uint, secret string) (string, error) {
	claims := &JwtCustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			Subject: username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
