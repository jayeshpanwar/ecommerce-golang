package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model

	Order_id   uint
	Product_id uint
	Quantity   uint
	Price      float32
}
