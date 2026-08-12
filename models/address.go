package models

import "gorm.io/gorm"

type Address struct {
	gorm.Model

	UserID uint

	FullName string `json:"full_name"`
	Phone    string `json:"phone"`

	AddressLine1 string `gorm:"column:address_line_1" json:"address_line_1"`
	AddressLine2 string `gorm:"column:address_line_2" json:"address_line_2"`
	City         string `json:"city"`
	State        string `json:"state"`
	Pincode      string `json:"pincode"`

	IsDefault bool `json:"is_default"`

	User User `json:"-"`
}
