package controllers

import (
	"ecommerce/dto"
	"ecommerce/services"

	"github.com/gin-gonic/gin"
)

func CreateReview(c *gin.Context) {

	var request dto.CreateReviewRequest

	if err := c.ShouldBindJSON(&request); err != nil {

		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	userID := c.MustGet("userID").(uint)

	review, err := services.CreateReview(userID, request)

	if err != nil {

		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	response := dto.BuildReviewResponse(*review)

	c.JSON(201, gin.H{

		"message": "Review created successfully",

		"review": response,
	})
}
