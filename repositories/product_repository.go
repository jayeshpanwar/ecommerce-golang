package repositories

import (
	"ecommerce/config"
	"ecommerce/models"
)

func CreateProduct(product *models.Product) error {
	return config.DB.Create(product).Error
}

func GetApprovedProducts(limit int, offset int) ([]models.Product, int64, error) {

	var total int64

	config.DB.
		Model(&models.Product{}).
		Where("status = ?", "approved").
		Count(&total)

	var products []models.Product

	err := config.DB.
		Where("status = ?", "approved").
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	return products, total, err
}

// --------------------------------------------------------------------------------------

func GetApprovedProductByID(id uint) (models.Product, error) {

	var product models.Product

	err := config.DB.
		Where("id = ? AND status = ?", id, "approved").
		First(&product).Error

	return product, err
}

//--------------------------------------------------------------------------------------

func GetProductByID(id uint) (*models.Product, error) {

	var product models.Product

	err := config.DB.
		First(&product, id).
		Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

//-------------------------------------------------------------------------

func DeleteProduct(product *models.Product) error {

	return config.DB.
		Delete(product).
		Error
}

//---------------------------------------------------------------------------------------------

func SaveProduct(product *models.Product) error {
	return config.DB.Save(product).Error
}

// ---------------------------------------------------------------------------------------------
func GetProductsBySellerID(
	sellerID uint,
) ([]models.Product, error) {

	var products []models.Product

	err := config.DB.
		Where("seller_id = ?", sellerID).
		Find(&products).
		Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

//----------------------------------------------------------------------------------------------

func GetPendingProducts() ([]models.Product, error) {

	var products []models.Product

	err := config.DB.
		Where("status = ?", "pending").
		Find(&products).
		Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

//----------------------------------------------------------------------------------------------
