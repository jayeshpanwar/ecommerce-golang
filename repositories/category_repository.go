package repositories

import (
	"ecommerce/config"
	"ecommerce/models"
)

func CreateCategory(category *models.Category) error {
	return config.DB.Create(category).Error
}
