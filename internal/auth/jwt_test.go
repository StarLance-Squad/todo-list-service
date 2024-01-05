package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateToken_EmptySecret(t *testing.T) {
	// Arrange
	username := "testuser"
	userID := uint(1)
	secret := ""

	// Act
	token, err := GenerateToken(username, userID, secret)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestGenerateToken_ValidInput(t *testing.T) {
	username := "testuser"
	userID := uint(123)
	secret := "testsecret"

	token, err := GenerateToken(username, userID, secret)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateToken_NonEmptyToken(t *testing.T) {
	username := "testuser"
	userID := uint(123)
	secret := "testsecret"

	token, err := GenerateToken(username, userID, secret)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateToken_ValidSubject(t *testing.T) {
	username := "testuser"
	userID := uint(123)
	secret := "testsecret"

	token, err := GenerateToken(username, userID, secret)

	assert.NoError(t, err)

	parsedToken, err := jwt.ParseWithClaims(token, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	assert.NoError(t, err)
	claims, ok := parsedToken.Claims.(*JwtCustomClaims)
	assert.True(t, ok)
	assert.Equal(t, username, claims.Subject)
}

func TestGenerateToken_EmptyUsername(t *testing.T) {
	username := ""
	userID := uint(123)
	secret := "testsecret"

	token, err := GenerateToken(username, userID, secret)

	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestGenerateToken_UserIDBoundary(t *testing.T) {
	username := "testuser"
	userID := uint(0) // or the max value of uint
	secret := "testsecret"

	token, err := GenerateToken(username, userID, secret)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
