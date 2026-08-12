package repositories

import (
	"ecommerce/config"
	"ecommerce/models"
)

func CreateCategory(category *models.Category) error {
	return config.DB.Create(category).Error
}

func GetCategories() ([]models.Category, error) {
	var categories []models.Category

	err := config.DB.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}
