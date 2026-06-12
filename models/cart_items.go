package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model

	CartID    uint
	ProductID uint
	Quantity  uint
}
