package dto

// LoginRequest 定义了用户登录时需要提供的请求体
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"` // 用户名
	Password string `json:"password" binding:"required" example:"123456"` // 密码
}

// LoginResponse 定义了成功登录后返回的数据结构
type LoginResponse struct {
	Token string `json:"token"` // JWT 令牌
}