package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"default:user"`
	IsActive bool   `gorm:"default:true"`

	Phone   string
	Address string
	City    string
	State   string
	Pincode string
}
