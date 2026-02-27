package dto

// LoginRequest 定义了用户登录时需要提供的请求体.
type LoginRequest struct {
	Username string `form:"username" binding:"required" example:"admin"`  // 用户名
	Password string `form:"password" binding:"required" example:"123456"` // 密码
}

// LoginResponse 定义了成功登录后返回的数据结构.
type LoginResponse struct {
	AccessToken   string   `json:"accessToken"`              // JWT 访问令牌
	MFARequired   bool     `json:"mfa_required,omitempty"`   // 是否需要 MFA
	RequiredTypes []string `json:"required_types,omitempty"` // 需要的 MFA 类型列表
}

// RefreshTokenResponse 定义了刷新令牌后返回的数据结构.
type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"` // 新的 JWT 访问令牌
}

// MFAVerifyRequest 定义了 MFA 验证时需要提供的请求体.
type MFAVerifyRequest struct {
	Code string `json:"code" binding:"required" example:"123456"` // MFA 验证码
}
