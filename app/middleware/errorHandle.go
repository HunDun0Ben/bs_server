package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/HunDun0Ben/bs_server/common/errors"
)

func WebErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// Check if there are errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if appErr, ok := err.(*errors.AppError); ok {
				c.JSON(appErr.Code, gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				})
				return
			}
			// Handle unknown errors
			c.JSON(500, gin.H{
				"code":    500,
				"message": "Internal Server Error",
			})
		}
	}
}
