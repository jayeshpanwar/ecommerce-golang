package controllers

import (
	"ecommerce/config"
	"ecommerce/models"
	"ecommerce/utils"

	"github.com/gin-gonic/gin"
)

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func Signup(c *gin.Context) {

	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid Credentials",
		})
		return
	}

	if req.Name == "" ||
		req.Email == "" ||
		req.Password == "" ||
		req.Role == "" {
		c.JSON(400, gin.H{
			"message": "All fields are required",
		})
		return
	}
	if req.Role != "user" &&
		req.Role != "seller" {
		c.JSON(400, gin.H{
			"Message": "Invalid role",
		})
		return
	}

	var existinguser models.User

	if err := config.DB.Where("email=?", req.Email).First(&existinguser).Error; err == nil {
		c.JSON(401, gin.H{
			"message": "Email already exists",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "Error Hasing Password",
		})
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(201, gin.H{
		"message": "User created successfully",
	})
}

//-----------------------------------LOGIN--------------------------------------------------------

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(404, gin.H{
			"erorr": err.Error(),
		})
		return
	}

	var user models.User

	if err := config.DB.Where("email=?", req.Email).First(&user).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid credentials",
		})
		return
	}

	err := utils.CheckPassword(user.Password, req.Password)

	if err != nil {
		c.JSON(401, gin.H{
			"message": "Invalid credentials",
		})
		return
	}

	token, err := utils.GenerateToken(
		user.ID,
		user.Role,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to generate token",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Login successful",
		"token":   token,
	})

}
