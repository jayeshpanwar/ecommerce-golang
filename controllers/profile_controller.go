package controllers

import ("ecommerce/models"
	"ecommerce/config"

	"github.com/gin-gonic/gin")

func updateProfile(c *gin.Context){

	userID:=c.MustGet("userID")

	var user map[string]interface{}

	if err:=c.ShouldBindJSON(&user).Error;err!=nil{
		c.JSON(400,gin.H{
			"message":"Invalid Data",
		})
	}

	if err:= config.DB.Find(&user,userID)


}