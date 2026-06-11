package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model

	User_id int
}
