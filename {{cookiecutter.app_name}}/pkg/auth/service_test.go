package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	signingKey := "test_signing_key"
	authService := NewAuthService(signingKey)

	userID := 1
	username := "testuser"

	tokenString, err := authService.CreateToken(userID, username)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse the token to verify its claims
	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(*Claim)
	assert.True(t, ok)
	assert.Equal(t, userID, claims.ID)
	assert.Equal(t, username, claims.Username)
	assert.WithinDuration(t, time.Now().Add(time.Hour*12), claims.ExpiresAt.Time, time.Minute)
}
func TestParseClaims(t *testing.T) {
	signingKey := "test_signing_key"
	authService := NewAuthService(signingKey)

	userID := 1
	username := "testuser"

	// Create a token to test ParseClaims
	tokenString, err := authService.CreateToken(userID, username)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Test valid token
	claims, err := authService.ParseClaims(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.ID)
	assert.Equal(t, username, claims.Username)

	// Test invalid token
	invalidTokenString := tokenString + "invalid"
	claims, err = authService.ParseClaims(invalidTokenString)
	assert.Error(t, err)
	assert.Nil(t, claims)

	// Test expired token
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claim{
		ID:       userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	})
	expiredTokenString, err := expiredToken.SignedString([]byte(signingKey))
	assert.NoError(t, err)

	claims, err = authService.ParseClaims(expiredTokenString)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

