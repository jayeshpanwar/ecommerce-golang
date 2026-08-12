package dto

import "ecommerce/models"

func BuildReviewResponse(review models.Review) ReviewResponse {

	return ReviewResponse{
		ID: review.ID,
		//ReviewerName: review.User.FullName,
		Rating:    review.Rating,
		Comment:   review.Comment,
		Images:    []string{},
		CreatedAt: review.CreatedAt,
	}
}
