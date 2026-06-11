package main

import (
	"ecommerce/config"
	"ecommerce/models"
	"ecommerce/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Cart{},
		&models.CartItem{},
	)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run(":8080")
}
