package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model

	UserID      uint
	TotalAmount float64
	Status      string

	OrderItems []OrderItem

	User User
}
