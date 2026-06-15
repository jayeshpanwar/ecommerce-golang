package routes

import (
	"ecommerce/controllers"
	"ecommerce/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.POST("/products", controllers.CreateProduct)
	router.GET("/products", controllers.GetProducts)
	router.GET("/products/:id", controllers.GetProductById)
	router.PUT("/products/:id", controllers.UpdateProduct)
	router.DELETE("/products/:id", controllers.DeleteProduct)
	router.PUT("/products/:id/approve", controllers.ApproveProduct)
	router.PUT("/products/:id/reject", controllers.RejectProduct)
	router.GET("/shop/products", controllers.GetApprovedProducts)
	router.POST("/cart/add", controllers.AddToCart)
	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.GET(
		"/profile",
		middlewares.JWTMiddleware(),
		func(c *gin.Context) {

			userID, _ := c.Get("userID")
			role, _ := c.Get("role")

			c.JSON(200, gin.H{
				"user_id": userID,
				"role":    role,
			})
		},
	)
}
