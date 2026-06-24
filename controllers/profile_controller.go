package controllers

import (
	"ecommerce/config"
	"ecommerce/models"

	"github.com/gin-gonic/gin"
)

type UpdateProfileRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Pincode string `json:"pincode"`
}

func GetProfile(c *gin.Context) {

	userID := c.MustGet("userID").(uint)

	var user models.User

	if err := config.DB.First(&user, userID).Error; err != nil {

		c.JSON(404, gin.H{
			"message": "User not found",
		})
		return
	}

	c.JSON(200, gin.H{})
}
