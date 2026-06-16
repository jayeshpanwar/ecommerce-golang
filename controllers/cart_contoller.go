package controllers

import (
	"ecommerce/config"
	"ecommerce/models"

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

//---------------------------DELETE SPECIFIC ITEM FROM CART--------------------------------------------------------

// ---------------------------VIEW CART--------------------------------------------------------
func ViewCart(c *gin.Context) {

	userID := c.MustGet("userID").(uint)

	var cart models.Cart

	if err := config.DB.Where("user_id=?", userID).First(&cart).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "Cart not found",
		})
		return
	}

	var cartItems []models.CartItem

	if err := config.DB.
		Where("cart_id = ?", cart.ID).
		Find(&cartItems).Error; err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})

		return
	}

	type CartItemResponse struct {
		ProductID   uint    `json:"product_id"`
		ProductName string  `json:"product_name"`
		Price       float64 `json:"price"`
		Quantity    uint    `json:"quantity"`
	}

	var response []CartItemResponse

	for _, item := range cartItems {

		var product models.Product

		if err := config.DB.
			First(&product, item.ProductID).
			Error; err != nil {

			continue
		}

		response = append(
			response,
			CartItemResponse{
				ProductID:   product.ID,
				ProductName: product.Name,
				Price:       product.Price,
				Quantity:    item.Quantity,
			},
		)
	}

	c.JSON(200, gin.H{
		"cart_id": cart.ID,
		"items":   response,
	})
}
