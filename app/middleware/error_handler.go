package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/HunDun0Ben/bs_server/app/pkg/bsvo"
)

func WebErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// 内部返回后
		// 查看 cxt 中是有错误.
		if len(c.Errors) > 0 {
			code := http.StatusInternalServerError
			message := http.StatusText(http.StatusInternalServerError)
			var data interface{}
			err := c.Errors.Last().Err
			if appErr, ok := err.(*bsvo.AppError); ok {
				code = appErr.Code
				message = appErr.Message
				data = appErr.Data
			}
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
				"data":    data,
			})
			c.Abort()
		}
	}
}
