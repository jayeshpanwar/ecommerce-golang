package services

import (
	"ecommerce/dto"
	"ecommerce/models"
	"ecommerce/repositories"
	"errors"
)

func CreateReview(
	userID uint,
	request dto.CreateReviewRequest,
) (*models.Review, error) {

	orderItem, err := repositories.GetOrderItemByID(request.OrderItemID)
	if err != nil {
		return nil, errors.New("order item not found")
	}

	order, err := repositories.GetOrderByID(orderItem.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.UserID != userID {
		return nil, errors.New("you cannot review this order")
	}

	if orderItem.Status != "delivered" {
		return nil, errors.New("product is not delivered yet")
	}

	if request.Rating < 1 || request.Rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	existingReview, err := repositories.GetReviewByOrderItemID(request.OrderItemID)
	if err != nil {
		return nil, err
	}

	if existingReview != nil {
		return nil, errors.New("review already exists")
	}

	review := models.Review{
		UserID:      userID,
		ProductID:   orderItem.ProductID,
		OrderItemID: orderItem.ID,
		Rating:      request.Rating,
		Comment:     request.Comment,
	}

	err = repositories.CreateReview(&review)
	if err != nil {
		return nil, err
	}

	return &review, nil
}
