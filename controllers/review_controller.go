package controllers

// import (
// 	"ecommerce/config"

// 	"github.com/gin-gonic/gin"
// )

// func CreateReview(c *gin.Context) {

// 	userID := c.MustGet("userID").(uint)

// 	var req CreateReviewRequest
// 	_ = CreateReviewRequest{}

// 	if err := c.ShouldBind(&req); err != nil {
// 		c.JSON(400, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if req.Rating < 1 || req.Rating > 5 {
// 		c.JSON(400, gin.H{
// 			"message": "Rating must be between 1 and 5",
// 		})
// 		return
// 	}

// 	var count int64

// 	config.DB.
// 		Table("orders").
// 		Joins("JOIN order_items ON order_items.order_id = orders.id").
// 		Where("orders.user_id = ? AND order_items.product_id = ? AND orders.status = ?",
// 			userID, req.ProductID, "delivered").
// 		Count(&count)

// 	if count == 0 {
// 		c.JSON(403, gin.H{
// 			"message": "You can only review purchased products",
// 		})
// 		return
// 	}

//}
