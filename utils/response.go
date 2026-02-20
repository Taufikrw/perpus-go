package utils

import "github.com/gin-gonic/gin"

func SendResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"status":  code,
		"message": message,
		"data":    data,
	})
}

func SendErrorResponse(c *gin.Context, code int, message string, errors interface{}) {
	if errors == nil {
		c.JSON(code, gin.H{
			"status":  code,
			"message": message,
		})
		return
	}

	c.JSON(code, gin.H{
		"status":  code,
		"message": message,
		"errors":  errors,
	})
}
