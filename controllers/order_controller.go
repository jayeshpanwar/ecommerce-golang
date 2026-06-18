package controllers

import (
	"ecommerce/config"
	"ecommerce/models"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Checkout(c *gin.Context) {

	userID := c.MustGet("userID").(uint)

	var cart models.Cart

	if err := config.DB.
		Where("user_id = ?", userID).
		First(&cart).Error; err != nil {

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

	type CheckoutItem struct {
		ProductID   uint    `json:"product_id"`
		ProductName string  `json:"product_name"`
		Price       float64 `json:"price"`
		Quantity    uint    `json:"quantity"`
		Subtotal    float64 `json:"subtotal"`
	}

	var items []CheckoutItem
	var total float64

	for _, item := range cartItems {

		var product models.Product

		if err := config.DB.
			First(&product, item.ProductID).Error; err != nil {

			continue
		}

		subtotal := product.Price * float64(item.Quantity)

		total += subtotal

		items = append(items, CheckoutItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	c.JSON(200, gin.H{
		"items":        items,
		"total_amount": total,
	})
}

func CreateOrder(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var cart models.Cart
		if err := tx.Where("user_id = ?", userID).First(&cart).Error; err != nil {
			c.JSON(400, gin.H{"message": "Cart not found"})
			return err
		}

		var cartItems []models.CartItem
		if err := tx.Where("cart_id = ?", cart.ID).Find(&cartItems).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return err
		}
		if len(cartItems) == 0 {
			c.JSON(400, gin.H{"message": "Cart is empty"})
			return fmt.Errorf("cart empty")
		}

		total := 0.0
		for _, item := range cartItems {
			var product models.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				c.JSON(500, gin.H{"error": "Product not found"})
				return err
			}
			if item.Quantity > uint(product.Stock) {
				c.JSON(400, gin.H{
					"error": fmt.Sprintf(
						"%s: only %d items available in stock",
						product.Name,
						product.Stock,
					),
				})
				return fmt.Errorf("insufficient stock")
			}
			total += product.Price * float64(item.Quantity)
		}

		order := models.Order{UserID: userID, TotalAmount: total, Status: "pending"}
		if err := tx.Create(&order).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return err
		}

		for _, item := range cartItems {
			var product models.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				c.JSON(500, gin.H{"error": "Product not found"})
				return err
			}

			orderItem := models.OrderItem{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				SellerID:  product.SellerID,
				Quantity:  item.Quantity,
				Price:     product.Price,
				Status:    "pending",
				Product:   product,
			}

			if err := tx.Create(&orderItem).Error; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return err
			}

			product.Stock -= int(item.Quantity)

			if err := tx.Save(&product).Error; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return err
			}
		}

		// Remove cart completely (items + cart row)
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return err
		}
		if err := tx.Delete(&cart).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return err
		}

		c.JSON(201, gin.H{
			"message":  "Order placed successfully",
			"order_id": order.ID,
			"total":    order.TotalAmount,
		})
		return nil
	})

	if err != nil {
		// Errors are already sent inside the transaction in most cases.
		return
	}
} //-----------------------------View Orders-----------------------------------------

func ViewOrders(c *gin.Context) {

	userID := c.MustGet("userID").(uint)

	type OrderItemResponse struct {
		ProductID   uint    `json:"product_id"`
		ProductName string  `json:"product_name"`
		Price       float64 `json:"price"`
		Quantity    uint    `json:"quantity"`
		Status      string  `json:"status"`
	}

	type OrderResponse struct {
		OrderID     uint                `json:"order_id"`
		OrderedAt   time.Time           `json:"ordered_at"`
		TotalAmount float64             `json:"total_amount"`
		Status      string              `json:"status"`
		Items       []OrderItemResponse `json:"items"`
	}

	var orders []models.Order

	if err := config.DB.Where("user_id = ?", userID).
		Order("created_at DESC").Preload("OrderItems.Product").
		Find(&orders).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var response []OrderResponse

	for _, order := range orders {

		var items []OrderItemResponse

		for _, item := range order.OrderItems {

			items = append(items, OrderItemResponse{
				ProductID:   item.ProductID,
				ProductName: item.Product.Name,
				Price:       item.Price,
				Quantity:    item.Quantity,
				Status:      item.Status,
			})
		}

		response = append(response, OrderResponse{
			OrderID:     order.ID,
			OrderedAt:   order.CreatedAt,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			Items:       items,
		})
	}
	c.JSON(200, gin.H{
		"data": response,
	})

}

//-----------------------------Seller View Orders-----------------------------------------

func GetSellerOrders(c *gin.Context) {

	sellerID := c.MustGet("userID").(uint)

	var orderItems []models.OrderItem

	if err := config.DB.
		Where("seller_id = ?", sellerID).
		Order("created_at DESC").
		Preload("Product").
		Preload("Order.User").
		Find(&orderItems).Error; err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	type SellerOrderResponse struct {
		OrderID     uint      `json:"order_id"`
		OrderItemID uint      `json:"order_item_id"`
		OrderedAt   time.Time `json:"ordered_at"`

		ProductID   uint   `json:"product_id"`
		ProductName string `json:"product_name"`

		Quantity uint    `json:"quantity"`
		Price    float64 `json:"price"`

		CustomerName  string `json:"customer_name"`
		CustomerEmail string `json:"customer_email"`

		OrderStatus string `json:"order_status"`
	}

	var response []SellerOrderResponse

	for _, item := range orderItems {

		response = append(response, SellerOrderResponse{
			OrderID:     item.Order.ID,
			OrderItemID: item.ID,
			OrderedAt:   item.Order.CreatedAt,

			ProductID:   item.Product.ID,
			ProductName: item.Product.Name,

			Quantity: item.Quantity,
			Price:    item.Price,

			CustomerName:  item.Order.User.Name,
			CustomerEmail: item.Order.User.Email,

			OrderStatus: item.Status,
		})
	}

	c.JSON(200, gin.H{
		"data": response,
	})
}

//-----------------------------Seller Update Order Status-----------------------------------------

func AcceptOrder(c *gin.Context) {

	orderItemID := c.Param("id")
	sellerID := c.MustGet("userID").(uint)

	var item models.OrderItem

	if err := config.DB.
		Where("id = ? AND seller_id = ?", orderItemID, sellerID).
		First(&item).Error; err != nil {

		c.JSON(404, gin.H{
			"message": "Order not found",
		})
		return
	}

	if item.Status != "pending" {
		c.JSON(400, gin.H{
			"message": "Order already processed",
		})
		return
	}

	item.Status = "accepted"

	if err := config.DB.Save(&item).Error; err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Order accepted",
	})
}

/////-------------REJECT ORDER----------------

func RejectOrder(c *gin.Context) {

	orderItemID := c.Param("id")
	sellerID := c.MustGet("userID").(uint)

	var item models.OrderItem

	if err := config.DB.
		Where("id = ? AND seller_id = ?", orderItemID, sellerID).
		First(&item).Error; err != nil {

		c.JSON(404, gin.H{
			"message": "Order not found",
		})
		return
	}
	if item.Status != "pending" {
		c.JSON(400, gin.H{
			"message": "Order already processed",
		})
		return
	}
	item.Status = "rejected"

	config.DB.Save(&item)

	c.JSON(200, gin.H{
		"message": "Order rejected",
	})
}

//-------Shipped Order-------------------

func ShipOrder(c *gin.Context) {

	orderItemID := c.Param("id")
	sellerID := c.MustGet("userID").(uint)

	var item models.OrderItem

	if err := config.DB.
		Where("id = ? AND seller_id = ?", orderItemID, sellerID).
		First(&item).Error; err != nil {

		c.JSON(404, gin.H{
			"message": "Order not found",
		})
		return
	}

	if item.Status != "accepted" {

		c.JSON(400, gin.H{
			"message": "Only accepted orders can be shipped",
		})
		return
	}

	item.Status = "shipped"

	if err := config.DB.Save(&item).Error; err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Order shipped",
	})
}

// ---------------- Delivered Order------------------

func DeliverOrder(c *gin.Context) {

	orderItemID := c.Param("id")
	sellerID := c.MustGet("userID").(uint)

	var item models.OrderItem

	if err := config.DB.
		Where("id = ? AND seller_id = ?", orderItemID, sellerID).
		First(&item).Error; err != nil {

		c.JSON(404, gin.H{
			"message": "Order not found",
		})
		return
	}

	if item.Status != "shipped" {

		c.JSON(400, gin.H{
			"message": "Only shipped orders can be delivered",
		})
		return
	}

	item.Status = "delivered"

	config.DB.Save(&item)

	c.JSON(200, gin.H{
		"message": "Order delivered",
	})
}
