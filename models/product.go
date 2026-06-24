package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model

	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock"`
	CategoryID  uint
	Category    Category `json:"category"`

	SellerID uint   `json:"seller_id"`
	Status   string `json:"status"`
}
