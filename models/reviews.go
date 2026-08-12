package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model

	UserID uint `json:"user_id"`

	ProductID uint `json:"product_id"`

	OrderItemID uint `json:"order_item_id"`

	Rating uint8 `json:"rating"`

	Comment string `json:"comment"`

	User User `gorm:"foreignKey:UserID"`

	Product Product `gorm:"foreignKey:ProductID"`

	OrderItem OrderItem `gorm:"foreignKey:OrderItemID"`

	Images []ReviewImage `gorm:"foreignKey:ReviewID"`
}
