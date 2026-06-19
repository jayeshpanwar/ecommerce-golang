package middlewares

import (
	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRole string) gin.HandlerFunc {

	return func(c *gin.Context) {

		role, exists := c.Get("role")

		if !exists {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		if role != allowedRole {
			c.JSON(403, gin.H{
				"message": "Forbidden",
			})
			c.Abort()
			return
		}
		c.Next()
	}

}
