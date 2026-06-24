package services

import (
	"errors"

	"ecommerce/models"
	"ecommerce/repositories"
	"ecommerce/utils"
)

func CreateProduct(product *models.Product) error {

	if product.Name == "" {
		return errors.New("product name is required")
	}

	if product.Price <= 0 {
		return errors.New("price must be greater than zero")
	}

	if product.Stock < 0 {
		return errors.New("invalid stock")
	}

	return repositories.CreateProduct(product)
}

//--------------------------------------------------------------------------------------------

func GetApprovedProducts(page int, limit int) ([]utils.ProductResponse, int64, error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	products, total, err := repositories.GetApprovedProducts(limit, offset)

	if err != nil {
		return nil, 0, err
	}

	var response []utils.ProductResponse

	for _, product := range products {

		response = append(
			response,
			utils.BuildProductResponse(product),
		)
	}

	return response, total, nil
}

// --------------------------------------------------------------------------------------------
func GetApprovedProductByID(id uint) (*utils.ProductResponse, error) {

	product, err :=
		repositories.GetApprovedProductByID(id)

	if err != nil {
		return nil, err
	}

	response :=
		utils.BuildProductResponse(product)

	return &response, nil
}

// -----------------------------------------------------------------------------------------------
var ErrProductNotFound = errors.New("product not found")
var ErrNotOwner = errors.New("you are not the owner of this product")

func DeleteProduct(productID uint, userID uint) error {

	product, err := repositories.GetProductByID(productID)

	if err != nil {
		return ErrProductNotFound
	}

	if product.SellerID != userID {
		return ErrNotOwner
	}

	return repositories.DeleteProduct(product)
}

//-----------------------------------------------------------------------------------------------

func ApproveProduct(productID uint) error {

	product, err :=
		repositories.GetProductByID(productID)

	if err != nil {
		return ErrProductNotFound
	}

	product.Status = "approved"

	return repositories.SaveProduct(product)
}

//-----------------------------------------------------------------------------------------------

func RejectProduct(productID uint) error {

	product, err := repositories.GetProductByID(productID)
	if err != nil {
		return ErrProductNotFound
	}

	if product.Status == "approved" {
		return errors.New("cannot reject an approved product")
	}

	product.Status = "rejected"
	return repositories.SaveProduct(product)
}

// -----------------------------------------------------------------------------------------------
func GetSellerProducts(
	sellerID uint,
) ([]utils.ProductResponse, error) {

	products, err :=
		repositories.GetProductsBySellerID(
			sellerID,
		)

	if err != nil {
		return nil, err
	}

	var response []utils.ProductResponse

	for _, product := range products {

		response = append(
			response,
			utils.BuildProductResponse(product),
		)
	}

	return response, nil
}

// -----------------------------------------------------------------------------------------------
func GetPendingProducts() ([]utils.ProductResponse, error) {

	products, err :=
		repositories.GetPendingProducts()

	if err != nil {
		return nil, err
	}

	var response []utils.ProductResponse

	for _, product := range products {

		response = append(
			response,
			utils.BuildProductResponse(product),
		)
	}

	return response, nil
}
