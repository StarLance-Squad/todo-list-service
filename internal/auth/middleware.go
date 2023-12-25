package auth

import (
	"github.com/labstack/echo/v4/middleware"
)

func JWTMiddlewareConfig(secretKey string) middleware.JWTConfig {
	return middleware.JWTConfig{
		Claims:     &JwtCustomClaims{},
		SigningKey: []byte(secretKey),
	}
}
