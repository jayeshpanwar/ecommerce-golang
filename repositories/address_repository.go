package repositories

import (
	"ecommerce/config"
	"ecommerce/models"
)

func GetAddresses(userID uint) ([]models.Address, error) {

	var addresses []models.Address
	err := config.DB.
		Where("user_id = ?", userID).
		Find(&addresses).Error

	return addresses, err
}
