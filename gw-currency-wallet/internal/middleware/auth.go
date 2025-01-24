package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwttoken "github.com/mizmorr/gw_currency/gw-currency-wallet/pkg/jwtToken"
)

func JWTAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := parts[1]

		err := jwttoken.Validate(tokenString, []byte(secretKey))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}
		userID, err := jwttoken.GetUserID(tokenString, []byte(secretKey))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}
		c.Set("user_id", userID)

		c.Next()
	}
}
