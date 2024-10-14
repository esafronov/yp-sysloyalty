package domain

import (
	"github.com/golang-jwt/jwt/v4"
)

type JwtCustomClaims struct {
	UserID int64
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	UserID int64
	jwt.RegisteredClaims
}
