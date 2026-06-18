package controllers

import (
	"ecommerce/config"
	"ecommerce/models"
	"ecommerce/utils"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {

	name := c.PostForm("name")
	description := c.PostForm("description")

	price, err := strconv.ParseFloat(
		c.PostForm("price"),
		64,
	)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid price",
		})
		return
	}

	stock, err := strconv.Atoi(
		c.PostForm("stock"),
	)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid stock",
		})
		return
	}
	userID := c.MustGet("userID").(uint)

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Image is required",
		})
		return
	}

	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !allowedTypes[ext] {
		c.JSON(400, gin.H{
			"error": "Only JPG, JPEG and PNG files are allowed",
		})
		return
	}

	if file.Size > 5*1024*1024 {
		c.JSON(400, gin.H{
			"error": "Image size must be less than 5 MB",
		})
		return
	}
	fileName := fmt.Sprintf(
		"%d_%s",
		time.Now().Unix(),
		file.Filename,
	)

	filePath := "./uploads/" + fileName

	if err := c.SaveUploadedFile(
		file,
		filePath,
	); err != nil {

		c.JSON(500, gin.H{
			"error": "Failed to save image",
		})
		return
	}

	product := models.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		SellerID:    userID,
		Status:      "pending",
		ImageURL:    "/uploads/" + fileName,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "Product Added Successfully",
		"product": product,
	})

}

//---------------------------------------------------------------------------------

func GetProducts(c *gin.Context) {

	var products []models.Product

	if err := config.DB.Where("status = ?", "approved").Find(&products).Error; err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var response []utils.ProductResponse

	for _, product := range products {

		response = append(
			response,
			utils.BuildProductResponse(
				product,
			),
		)
	}

	c.JSON(200, gin.H{
		"message": "Products fetched successfully",
		"data":    response,
	})

}

//----------------------------------GetProductByID-----------------------------------------

func GetProductById(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid product ID",
		})
		return
	}

	var product models.Product

	// Only expose approved products to users.
	if err := config.DB.Where("id = ? AND status = ?", id, "approved").First(&product).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Product Not Found",
		})
		return
	}

	response := utils.BuildProductResponse(product)

	c.JSON(200, gin.H{
		"message": "Product fetched by ID",
		"data":    response,
	})

}

//---------------------------------------------------------------------------------

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product

	//c.ShouldBindJSON(&product)    // The request body can only be read once.

	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Product Not Found",
		})
		return
	}

	userID := c.MustGet("userID").(uint)

	if product.SellerID != userID {
		c.JSON(403, gin.H{
			"message": "Forbidden: You are not the owner of this product",
		})
		return
	}

	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	delete(input, "seller_id")
	delete(input, "status")
	delete(input, "id")

	if err := config.DB.Model(&product).Updates(input).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Product Updated",
	})

}

//-------------------------------------------------------------------------------------

func DeleteProduct(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	var product models.Product

	err = config.DB.First(&product, id).Error
	if err != nil {
		c.JSON(404, gin.H{
			"error": "ID not exists",
		})
		return
	}

	userID := c.MustGet("userID").(uint)

	if product.SellerID != userID {
		c.JSON(403, gin.H{
			"message": "Forbidden: You are not the owner of this product",
		})
		return
	}

	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Product Deleted Successfully",
	})

}

//-------------------------------------------------------------------------------------

func ApproveProduct(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	var product models.Product

	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Product Not Found",
		})
		return
	}

	product.Status = "approved"

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Product Approved",
	})

}

//-------------------------------------------------------------------------------------

func RejectProduct(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	var product models.Product

	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Product Not Found",
		})
		return
	}

	product.Status = "rejected"

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Product Rejected",
	})

}

//-------------------------------------------------------------------------------------

func GetApprovedProducts(c *gin.Context) {

	var products []models.Product

	if err := config.DB.Where("status=?", "approved").Find(&products).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	var response []utils.ProductResponse

	for _, product := range products {
		response = append(
			response,
			utils.BuildProductResponse(product),
		)
	}

	c.JSON(200, gin.H{
		"data": response,
	})
}

// -----------------------GET SELLER PRODUCTS---------------------------------------------
func GetSellerProducts(c *gin.Context) {

	userID := c.MustGet("userID").(uint)

	var products []models.Product

	if err := config.DB.Where("seller_id=?", userID).Find(&products).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	var response []utils.ProductResponse

	for _, product := range products {
		response = append(
			response,
			utils.BuildProductResponse(product),
		)
	}

	c.JSON(200, gin.H{
		"data": response,
	})

}

// -----------------------GET PENDING PRODUCTS (ADMIN)---------------------------------
func GetPendingSellerProducts(c *gin.Context) {

	var products []models.Product

	if err := config.DB.Where("status = ?", "pending").Find(&products).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	var response []utils.ProductResponse

	for _, product := range products {
		response = append(
			response,
			utils.BuildProductResponse(product),
		)
	}

	c.JSON(200, gin.H{
		"data": response,
	})
}
