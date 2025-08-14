package dto

// AppRes 是应用统一的泛型响应结构体.
type AppRes[T any] struct {
	Code    int     `json:"code"`              // 业务码
	Message *string `json:"message,omitempty"` // 响应消息, nil 时忽略
	Data    T       `json:"data"`              // 响应数据
}

// 它仅用于 Swagger 文档生成，实际业务代码不使用.
type SwaggerResponse struct {
	Code    int         `json:"code"`    // 业务码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

func NewBaseRes[T any](code int, msg *string, data T) AppRes[T] {
	return AppRes[T]{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}
