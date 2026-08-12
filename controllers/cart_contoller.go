package controllers

import (
	"ecommerce/config"
	"ecommerce/models"
	"ecommerce/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddToCartRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

func AddToCart(c *gin.Context) {

	var req AddToCartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid Data",
		})
		return
	}

	//------------Check Product Approved or Not------------------------------------------
	var product models.Product

	if err := config.DB.Where("status=?", "approved").First(&product, req.ProductID).Error; // First(..., id) automatically adds an ID(PK) condition.
	err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid Product ID",
		})
		return
	}

	var cart models.Cart

	//------------Check Stock of Product------------------------------

	if req.Quantity > uint(product.Stock) {
		c.JSON(400, gin.H{
			"message": "Insufficient stock",
		})
		return
	}

	//------------Check Cart already available and create new if not------------------------------------------

	userID := c.MustGet("userID").(uint)

	if err := config.DB.Where("user_id=?", userID).First(&cart).Error; err != nil {
		cart = models.Cart{
			UserID: userID,
		}
		if err := config.DB.Create(&cart).Error; err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
	//------------Check Product Already in cart and inc. quantity------------------------------------------

	var cartItem models.CartItem

	result := config.DB.Where("cart_id=? AND product_id=?", cart.ID, req.ProductID).First(&cartItem)

	if result.Error == nil {

		newQuantity := cartItem.Quantity + req.Quantity

		if newQuantity > uint(product.Stock) { //check stock availability
			c.JSON(400, gin.H{
				"message": "Insufficient stock",
			})
			return
		}
		cartItem.Quantity = newQuantity

		config.DB.Save(&cartItem)
		c.JSON(200, gin.H{
			"message": "Cart Updated",
			"cart":    cartItem,
		})
		return
	}

	//-----------------------Add product in cart---------------------------------------------

	cartItem = models.CartItem{
		CartID:    cart.ID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := config.DB.Create(&cartItem).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Product added to cart",
		"cart":    cartItem,
	})

}

// ---------------------------DELETE SPECIFIC ITEM FROM CART--------------------------------------------------------
func DeleteFromCart(c *gin.Context) {

	cartItemID := c.Param("id")
	userID := c.MustGet("userID").(uint)

	err := services.DeleteFromCart(userID, cartItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart item deleted",
	})
}

// ---------------------------VIEW CART--------------------------------------------------------
func ViewCart(c *gin.Context) {

	userID := c.MustGet("userID").(uint)

	response, err := services.ViewCart(userID)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Cart not found",
		})
		return
	}

	c.JSON(200, response)
}

// /-------------------------DECREASE QUANTITY OF ITEM IN CART--------------------------------------------------------
func DecreaseCartItemQuantity(c *gin.Context) {

	cartItemID := c.Param("id")
	userID := c.MustGet("userID").(uint)

	var cart models.Cart

	if err := config.DB.Where("user_id=?", userID).First(&cart).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "Cart not found",
		})
		return
	}

	var cartItem models.CartItem

	fmt.Println("Cart Item ID:", cartItemID)
	fmt.Println("Cart ID:", cart.ID)

	if err := config.DB.Where("id=? AND cart_id=?", cartItemID, cart.ID).First(&cartItem).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "Cart item not found",
		})
		return
	}

	if cartItem.Quantity > 1 {
		cartItem.Quantity -= 1
		if err := config.DB.Save(&cartItem).Error; err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "Cart item quantity decreased",
			"cart":    cartItem,
		})
	} else if cartItem.Quantity == 1 {
		if err := config.DB.Delete(&cartItem).Error; err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "Cart item removed",
		})
	}
}
