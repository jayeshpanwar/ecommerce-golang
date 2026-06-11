package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model

	User_id      uint
	Total_amount uint
	Status       string
}
