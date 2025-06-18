package dto

type AppRes[T any] struct {
	Code    int     `json:"code"`              // 始终序列化
	Message *string `json:"message,omitempty"` // nil 时忽略
	Data    T       `json:"data"`
}

func NewBaseRes[T any](code int, msg *string, data T) AppRes[T] {
	return AppRes[T]{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}
