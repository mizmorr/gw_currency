package jwttoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokensOption struct {
	UserID        int64         `json:"user_id"`
	RefreshExp    time.Duration `json:"refresh_exp"`
	AccessExp     time.Duration `json:"access_exp"`
	SecretRefresh string        `json:"secret_refresh"`
	SecretAccess  string        `json:"secret_access"`
}

type CustomClaims struct {
	UserID int64
	jwt.RegisteredClaims
}
