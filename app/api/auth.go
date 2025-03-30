package api

import (
	"net/http"

	"github.com/HunDun0Ben/bs_server/app/middleware"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 这里应该添加实际的用户验证逻辑
	// 为了演示，我们只做简单的验证
	if req.Username == "admin" && req.Password == "password" {
		token, err := middleware.GenerateToken(req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
	}
}
