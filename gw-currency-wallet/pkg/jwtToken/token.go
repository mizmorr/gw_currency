package jwttoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(options *TokensOption) (string, string, error) {
	accessToken, err := generateToken(options.AccessExp, options.SecretAccess, options.UserID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateToken(options.RefreshExp, options.SecretRefresh, options.UserID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func generateToken(expTime time.Duration, secret string, id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   string(rune(id)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	return token.SignedString([]byte(secret))
}
