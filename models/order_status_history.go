package models

import "gorm.io/gorm"

type OrderStatusHistory struct {
	gorm.Model

	OrderID uint `json:"order_id"`

	Status string `json:"status"`

	UpdatedBy uint `json:"updated_by"`
}
