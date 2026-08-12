package services

import (
	"ecommerce/models"
	"ecommerce/repositories"
)

func GetAddresses(userID uint) ([]models.Address, error) {

	addresses, err := repositories.GetAddresses(userID)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}
