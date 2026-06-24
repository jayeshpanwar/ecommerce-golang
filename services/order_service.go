package services

import (
	"ecommerce/models"
	"ecommerce/repositories"
)

func AddOrderHistory(
	orderID uint,
	status string,
	updatedBy uint,
) error {

	history := models.OrderStatusHistory{
		OrderID:   orderID,
		Status:    status,
		UpdatedBy: updatedBy,
	}

	return repositories.CreateOrderStatusHistory(
		&history,
	)
}

func GetOrderHistory(
	orderID uint,
) ([]models.OrderStatusHistory, error) {

	return repositories.GetOrderHistory(orderID)
}
