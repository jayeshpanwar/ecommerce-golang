package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model

	Cart_id    uint
	Product_id uint
	Quantity   uint
}
