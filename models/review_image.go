package models

import "gorm.io/gorm"

type ReviewImage struct {
	gorm.Model

	ReviewID uint `json:"review_id"`

	ImageURL string `json:"image_url"`

	Review Review `gorm:"foreignKey:ReviewID"`
}
