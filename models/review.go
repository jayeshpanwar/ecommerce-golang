package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model

	UserID    uint
	ProductID uint
	OrderID   uint

	Rating  int // 1 to 5
	Comment string
}
