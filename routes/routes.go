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
		seller.GET("/products", controllers.GetSellerProducts)
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

	// Authenticated Routes (Role: user)
	user := router.Group("/")
	user.Use(
		middlewares.JWTMiddleware(),
		middlewares.RoleMiddleware("user"),
	)
	{
		// User sees only approved products.
		user.GET("/products", controllers.GetProducts)
		user.GET("/products/:id", controllers.GetProductById)
		user.POST("/cart/add", controllers.AddToCart)
		user.GET("/cart", controllers.ViewCart)
		//user.GET("/profile", controllers.Profile)
	}

	// Authenticated Routes (Role: seller) - for seller product management UI
	sellerUser := router.Group("/")
	sellerUser.Use(
		middlewares.JWTMiddleware(),
		middlewares.RoleMiddleware("seller"),
	)
	{
		// Seller can manage products via CRUD under /seller group.
		// No user-facing product listing needed here.
		// Keeping empty to allow future seller-specific endpoints.
	}

}
