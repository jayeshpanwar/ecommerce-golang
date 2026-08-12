package dto

type CreateOrderRequest struct {
	AddressID     uint   `json:"address_id"`
	PaymentMethod string `json:"payment_method"`
}
