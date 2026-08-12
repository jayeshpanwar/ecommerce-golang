package services

import (
	"ecommerce/dto"
	"ecommerce/repositories"
	"errors"
)

func DeleteFromCart(userID uint, cartItemID string) error {

	//  Get cart
	cart, err := repositories.GetCartByUserID(userID)
	if err != nil {
		return errors.New("cart not found")
	}

	//Get cart item
	item, err := repositories.GetCartItemByID(cart.ID, cartItemID)
	if err != nil {
		return errors.New("cart item not found")
	}

	//  Delete item
	if err := repositories.DeleteCartItem(&item); err != nil {
		return err
	}

	return nil
}

////--------------------------VIEW CART--------------------------------------------------------

func ViewCart(userID uint) (*dto.CartResponse, error) {

	cart, err := repositories.GetCartByUserID(userID)
	if err != nil {
		return nil, errors.New("cart not found")

	}

	return repositories.GetCartItems(cart.ID)
}
