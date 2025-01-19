package jwttoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Validate(tokenString string, jwtSecret []byte) error {
	tokenParsed, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidType
		}
		return jwtSecret, nil
	})
	if err != nil {
		return err
	}
	expTime, err := tokenParsed.Claims.GetExpirationTime()

	if err != nil || time.Now().After(expTime.Time) {
		return err
	}

	return nil
}
