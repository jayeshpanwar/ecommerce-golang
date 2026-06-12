package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model

	OrderID   uint
	ProductID uint
	Quantity  uint
	Price     float32
}
