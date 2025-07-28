package dto

// LoginRequest 定义了用户登录时需要提供的请求体
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"` // 用户名
	Password string `json:"password" binding:"required" example:"123456"` // 密码
}

// LoginResponse 定义了成功登录后返回的数据结构
type LoginResponse struct {
	AccessToken  string `json:"access_token"`  // JWT 访问令牌
	RefreshToken string `json:"refresh_token"` // JWT 刷新令牌
}

// RefreshTokenRequest 定义了刷新令牌的请求体
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"` // JWT 刷新令牌
}

// RefreshTokenResponse 定义了刷新令牌后返回的数据结构
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"` // 新的 JWT 访问令牌
}
