package utils

import (
	"ecommerce/models"
	"os"
)

type ProductResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Status      string  `json:"status"`
	ImageURL    string  `json:"image_url"`
}

func BuildProductResponse(
	product models.Product,
) ProductResponse {

	baseURL := os.Getenv("BASE_URL")

	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Status:      product.Status,
		ImageURL:    baseURL + product.ImageURL,
	}
}
