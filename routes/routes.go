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
		seller.GET("/categories", controllers.GetCategories)
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
		admin.POST("/categories", controllers.CreateCategory)

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
		user.POST("/payments/verify", controllers.VerifyPayment)
		user.GET("/orders/:id/payments", controllers.GetPaymentDetails)
		user.GET("/orders/:id/tracking", controllers.GetOrderTracking)
		user.POST("/payment/COD/:order_id", controllers.PlaceCODOrder)
		user.POST("/addresses", controllers.AddAddress)
		user.GET("/addresses", controllers.GetAddresses)
		user.GET("/products/category/:id", controllers.GetProductsByCategory)
		user.GET("/products/search", controllers.SearchProducts)
		user.POST("/reviews", controllers.CreateReview)
		//user.GET("/products/:id/reviews", controllers.GetProductReviews)

		// user.GET("/addresses/:id", controllers.GetAddressByID)
		// user.PUT("/addresses/:id", controllers.UpdateAddress)
		// user.DELETE("/addresses/:id", controllers.DeleteAddress)

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
