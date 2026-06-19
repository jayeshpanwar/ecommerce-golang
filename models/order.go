package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model

	UserID            uint
	TotalAmount       float64
	Status            string
	PaymentStatus     string // pending, paid, failed
	RazorpayOrderID   string // important
	RazorpayPaymentID string

	OrderItems []OrderItem

	User User
}
