package controllers

import (
	"ecommerce/models"
	"ecommerce/services"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {

	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := services.CreateCategory(&category); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "Category created successfully",
		"data":    category,
	})
}

func GetCategories(c *gin.Context) {

	categories, err := services.GetCategories()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch categories",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": categories,
	})
}
