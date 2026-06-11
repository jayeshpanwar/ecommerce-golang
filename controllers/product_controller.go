package controllers

import (
	"ecommerce/config"
	"ecommerce/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {

	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	product.Status = "pending"

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

	if err := config.DB.Find(&products).Error; err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Products fetched successfully",
		"data":    products,
	})

}

//---------------------------------------------------------------------------------

func GetProductById(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid product ID",
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

	c.JSON(200, gin.H{
		"message": "Product fetched by ID",
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

	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

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

	c.JSON(200, gin.H{
		"data": products,
	})
}
