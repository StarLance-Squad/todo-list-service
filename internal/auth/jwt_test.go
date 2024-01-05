package auth

import (
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
