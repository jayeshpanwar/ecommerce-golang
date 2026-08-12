package services

import (
	"ecommerce/models"
	"ecommerce/repositories"
)

func CreateCategory(category *models.Category) error {

	return repositories.CreateCategory(category)
}

func GetCategories() ([]models.Category, error) {
	return repositories.GetCategories()
}
