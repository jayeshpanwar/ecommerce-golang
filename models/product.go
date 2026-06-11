package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model

	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`

	SellerID uint   `json:"seller_id"`
	Status   string `json:"status"`
}
