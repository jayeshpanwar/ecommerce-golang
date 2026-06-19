package models

import "gorm.io/gorm"

type ReviewImage struct {
	gorm.Model

	ReviewID uint
	ImageURL string
}
