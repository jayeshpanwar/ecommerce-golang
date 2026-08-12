package controllers

import (
	"ecommerce/config"
	"ecommerce/dto"
	"ecommerce/models"
	"ecommerce/services"

	"github.com/gin-gonic/gin"
)

func AddAddress(c *gin.Context) {

	userID := c.MustGet("userID").(uint)

	var req dto.CreateAddressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error":   "validation failed",
			"details": err.Error(),
		})
		return
	}

	address := models.Address{
		UserID:       userID,
		FullName:     req.FullName,
		Phone:        req.Phone,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		Pincode:      req.Pincode,
	}

	if err := config.DB.Create(&address).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "Address added successfully",
		"data":    address,
	})
}

func GetAddresses(c *gin.Context) {

	userID := c.MustGet("userID").(uint)

	addresses, err := services.GetAddresses(userID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch addresses",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": addresses,
	})
}
