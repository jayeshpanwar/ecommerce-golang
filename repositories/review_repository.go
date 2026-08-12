package repositories

import (
	"ecommerce/config"
	"ecommerce/models"
	"errors"

	"gorm.io/gorm"
)

func GetOrderItemByID(id uint) (*models.OrderItem, error) {

	var orderItem models.OrderItem

	err := config.DB.
		First(&orderItem, id).Error

	if err != nil {
		return nil, err
	}

	return &orderItem, nil
}
func GetOrderByID(id uint) (*models.Order, error) {

	var order models.Order

	err := config.DB.
		First(&order, id).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func GetReviewByOrderItemID(orderItemID uint) (*models.Review, error) {

	var review models.Review

	err := config.DB.
		Where("order_item_id = ?", orderItemID).
		First(&review).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &review, nil
}

func CreateReview(review *models.Review) error {

	return config.DB.
		Create(review).Error
}
