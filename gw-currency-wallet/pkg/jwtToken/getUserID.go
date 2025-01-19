package jwttoken

import "github.com/golang-jwt/jwt/v5"

func GetUserID(tokenString string, jwtSecret []byte) (int64, error) {
	tokenParsed, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidType
		}
		return jwtSecret, nil
	})
	if err != nil {
		return 0, jwt.ErrSignatureInvalid
	}

	claims, ok := tokenParsed.Claims.(*CustomClaims)

	if !ok {
		return 0, jwt.ErrTokenInvalidClaims
	}
	return claims.UserID, nil
}
