package middlewares

import (
	"ecommerce/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader :=
			c.GetHeader("Authorization")

		if authHeader == "" {

			c.JSON(401, gin.H{
				"message": "Authorization header missing",
			})

			c.Abort()
			return
		}

		tokenString :=
			strings.TrimPrefix(
				authHeader,
				"Bearer ",
			)

		token, err :=
			utils.ValidateToken(
				tokenString,
			)

		if err != nil || !token.Valid {

			c.JSON(401, gin.H{
				"message": "Invalid token",
			})

			c.Abort()
			return
		}

		claims :=
			token.Claims.(jwt.MapClaims)

		c.Set(
			"userID",
			claims["user_id"],
		)

		c.Set(
			"role",
			claims["role"],
		)

		c.Next()
	}
}
