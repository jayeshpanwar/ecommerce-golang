package repositories

import (
	"ecommerce/config"
	"ecommerce/dto"
	"ecommerce/models"
)

func GetCartByUserID(userID uint) (models.Cart, error) {
	var cart models.Cart

	err := config.DB.
		Where("user_id = ?", userID).
		First(&cart).Error

	return cart, err
}

func GetCartItemByID(cartID uint, itemID string) (models.CartItem, error) {
	var item models.CartItem

	err := config.DB.
		Where("id = ? AND cart_id = ?", itemID, cartID).
		First(&item).Error

	return item, err
}

func DeleteCartItem(item *models.CartItem) error {
	return config.DB.Delete(item).Error
}

func GetCartItems(cartID uint) (*dto.CartResponse, error) {

	var items []dto.CartItemResponse

	err := config.DB.
		Table("cart_items").
		Select(`
			cart_items.id as cart_item_id,
			products.id as product_id,
			products.name as product_name,
			products.price,
			cart_items.quantity
		`).
		Joins("JOIN products ON products.id = cart_items.product_id").
		Where("cart_items.cart_id = ?", cartID).
		Where("cart_items.deleted_at IS NULL").
		Scan(&items).Error

	if err != nil {
		return nil, err
	}

	return &dto.CartResponse{
		CartID: cartID,
		Items:  items,
	}, nil
}
