package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService provides authentication services, including signing and verifying tokens.
// It uses a signing key to generate and validate authentication tokens.
type AuthService struct {
	signingKey string
}

// NewAuthService creates a new instance of AuthService with the provided signing key.
// signingKey: the key used for signing tokens.
// Returns a pointer to the newly created AuthService.
func NewAuthService(signingKey string) *AuthService {
	return &AuthService{signingKey: signingKey}
}

// CreateToken generates a JWT token for a given user ID and username.
// The token is signed using the AuthService's signing key and is valid for 12 hours.
//
// Parameters:
//   userID - the ID of the user for whom the token is being created
//   username - the username of the user for whom the token is being created
//
// Returns:
//   A signed JWT token as a string, or an error if the token could not be created.
func (a *AuthService) CreateToken(userID int, username string) (string, error) {
	claim := &Claim{
		ID:       userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(a.signingKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseClaims parses a JWT token string and returns the claims contained within it.
// It takes a token string as input and returns a pointer to the Claim struct and an error.
// If the token is invalid or the signature is invalid, it returns an error.
func (a *AuthService) ParseClaims(tokenString string) (*Claim, error) {
	claims := &Claim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.signingKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
