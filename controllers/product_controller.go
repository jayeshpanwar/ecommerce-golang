package controllers

import (
	"ecommerce/config"
	"ecommerce/models"
	"ecommerce/services"
	"ecommerce/utils"
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {

	name := c.PostForm("name")
	description := c.PostForm("description")
	categoryID, err := strconv.ParseUint(
		c.PostForm("category_id"),
		10,
		64,
	)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

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

	ext := strings.ToLower(
		filepath.Ext(file.Filename),
	)

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
		CategoryID:  uint(categoryID),
		ImageURL:    "/uploads/" + fileName,
	}

	if err := services.CreateProduct(&product); err != nil {

		c.JSON(400, gin.H{
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

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	response, total, err := services.GetApprovedProducts(page, limit)

	if err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	totalPages := int(math.Ceil(
		float64(total) / float64(limit),
	))

	c.JSON(200, gin.H{
		"page":          page,
		"limit":         limit,
		"totalProducts": total,
		"totalPages":    totalPages,
		"products":      response,
	})
}

// ----------------------------------GetProductByID-----------------------------------------
func GetProductById(c *gin.Context) {

	id, err :=
		strconv.ParseUint(
			c.Param("id"),
			10,
			64,
		)

	if err != nil {

		c.JSON(400, gin.H{
			"message": "Invalid product ID",
		})
		return
	}

	product, err :=
		services.GetApprovedProductByID(
			uint(id),
		)

	if err != nil {

		c.JSON(404, gin.H{
			"message": "Product not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Product fetched successfully",
		"data":    product,
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

	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	userID := c.MustGet("userID").(uint)

	err = services.DeleteProduct(
		uint(id),
		userID,
	)

	if err != nil {

		switch err {

		case services.ErrProductNotFound:
			c.JSON(404, gin.H{
				"error": err.Error(),
			})

		case services.ErrNotOwner:
			c.JSON(403, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(500, gin.H{
				"error": "Internal Server Error",
			})
		}

		return
	}

	c.JSON(200, gin.H{
		"message": "Product Deleted Successfully",
	})
}

//-------------------------------------------------------------------------------------

func ApproveProduct(c *gin.Context) {

	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {

		c.JSON(400, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	err = services.ApproveProduct(
		uint(id),
	)

	if err != nil {

		switch err {

		case services.ErrProductNotFound:
			c.JSON(404, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(500, gin.H{
				"error": "Internal Server Error",
			})
		}

		return
	}

	c.JSON(200, gin.H{
		"message": "Product Approved",
	})
}

// -------------------------------------------------------------------------------------
func RejectProduct(c *gin.Context) {

	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {

		c.JSON(400, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	err = services.RejectProduct(
		uint(id),
	)

	if err != nil {

		switch err {

		case services.ErrProductNotFound:
			c.JSON(404, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(500, gin.H{
				"error": "Internal Server Error",
			})
		}

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

	products, err :=
		services.GetSellerProducts(
			userID,
		)

	if err != nil {

		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": products,
	})
}

// -----------------------GET PENDING PRODUCTS (ADMIN)---------------------------------
func GetPendingSellerProducts(c *gin.Context) {

	products, err := services.GetPendingProducts()

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": products,
	})
}
