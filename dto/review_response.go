package dto

import "time"

type ReviewResponse struct {
	ID           uint      `json:"id"`
	ReviewerName string    `json:"reviewer_name"`
	Rating       uint8     `json:"rating"`
	Comment      string    `json:"comment"`
	Images       []string  `json:"images"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateReviewRequest struct {
	OrderItemID uint   `json:"order_item_id" binding:"required"`
	Rating      uint8  `json:"rating" binding:"required,gte=1,lte=5"`
	Comment     string `json:"comment"`
}
