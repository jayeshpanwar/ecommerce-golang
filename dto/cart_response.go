package dto

type CartItemResponse struct {
	CartItemID  uint    `json:"cart_item_id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    uint    `json:"quantity"`
}

type CartResponse struct {
	CartID uint               `json:"cart_id"`
	Items  []CartItemResponse `json:"items"`
}
