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
	router.POST("/webhooks/razorpay", controllers.RazorpayWebhook)

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
		seller.GET("/orders", controllers.GetSellerOrders)
		seller.PUT("/orders/:id/accept", controllers.AcceptOrder)
		seller.PUT("/orders/:id/reject", controllers.RejectOrder)
		seller.PUT("/orders/:id/ship", controllers.ShipOrder)
		seller.PUT("/orders/:id/deliver", controllers.DeliverOrder)
	}

	// Admin Routes
	admin := router.Group("/admin")
	admin.Use(
		middlewares.JWTMiddleware(),
		middlewares.RoleMiddleware("admin"),
	)
	{
		admin.GET("/products/pending", controllers.GetPendingSellerProducts)
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

		user.GET("/products/:id", controllers.GetProductById)
		user.POST("/cart/add", controllers.AddToCart)
		user.GET("/cart", controllers.ViewCart)
		user.PUT("/cart/items/:id/decrease", controllers.DecreaseCartItemQuantity)
		user.DELETE("/cart/items/:id", controllers.DeleteFromCart)
		user.GET("/checkout/summary", controllers.Checkout)
		user.POST("/orders", controllers.CreateOrder)
		user.GET("/orders", controllers.ViewOrders)
		user.POST("/payments/create-order/:order_id", controllers.CreatePaymentOrder)
		user.GET("/orders/:id/payments", controllers.GetPaymentDetails)
		//user.POST("/reviews", controllers.CreateReview)

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
