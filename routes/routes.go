package routes

import (
	"ecommerce/controllers"

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
}
