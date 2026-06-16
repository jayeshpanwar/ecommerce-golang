package routes

import (
	"ecommerce/controllers"
	"ecommerce/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	// Public Routes
	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)

	router.GET("/products", controllers.GetProducts)
	router.GET("/products/:id", controllers.GetProductById)
	router.GET("/shop/products", controllers.GetApprovedProducts)

	// Seller Routes
	seller := router.Group("/seller")
	seller.Use(
		middlewares.JWTMiddleware(),
		middlewares.RoleMiddleware("seller"),
	)
	{
		seller.POST("/products", controllers.CreateProduct)
		seller.PUT("/products/:id", controllers.UpdateProduct)
		seller.DELETE("/products/:id", controllers.DeleteProduct)
	}

	// Admin Routes
	admin := router.Group("/admin")
	admin.Use(
		middlewares.JWTMiddleware(),
		middlewares.RoleMiddleware("admin"),
	)
	{
		admin.PUT("/products/:id/approve", controllers.ApproveProduct)
		admin.PUT("/products/:id/reject", controllers.RejectProduct)
	}

	// Authenticated Routes
	auth := router.Group("/")
	auth.Use(middlewares.JWTMiddleware())
	{
		auth.POST("/cart/add", controllers.AddToCart)
		auth.GET("/cart", controllers.ViewCart)
		//auth.GET("/profile", controllers.Profile)
	}
}
