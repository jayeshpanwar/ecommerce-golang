package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model

	UserID      uint
	TotalAmount float64
	Status      string
	AddressID   uint
	Address     Address

	PaymentStatus     string // pending, cod, paid, failed
	RazorpayOrderID   string // important
	RazorpayPaymentID string
	PaymentMethod     string // cod, razorpay

	OrderItems []OrderItem

	User User
}
