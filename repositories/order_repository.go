package repositories

import (
	"ecommerce/config"
	"ecommerce/models"
)

func CreateOrderStatusHistory(
	history *models.OrderStatusHistory,
) error {

	return config.DB.Create(history).Error
}

func GetOrderHistory(
	orderID uint,
) ([]models.OrderStatusHistory, error) {

	var history []models.OrderStatusHistory

	err := config.DB.
		Where("order_id = ?", orderID).
		Order("created_at asc").
		Find(&history).
		Error

	return history, err
}
