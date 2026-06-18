package utils

import "github.com/gin-gonic/gin"

func GetImageURL(
	c *gin.Context,
	imagePath string,
) string {

	scheme := "http"

	if c.Request.TLS != nil {
		scheme = "https"
	}

	return scheme +
		"://" +
		c.Request.Host +
		imagePath
}
