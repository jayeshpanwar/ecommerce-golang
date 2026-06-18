package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model

	OrderID   uint
	ProductID uint
	Quantity  uint
	Price     float64

	Product  Product
	Order    Order
	SellerID uint
	Status   string
}
