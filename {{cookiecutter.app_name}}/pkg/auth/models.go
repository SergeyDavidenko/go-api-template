package auth

import "github.com/golang-jwt/jwt/v5"

type Claim struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
