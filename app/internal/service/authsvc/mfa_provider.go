package authsvc

import "context"

// MFAProvider 定义了所有二次验证方式必须实现的接口
type MFAProvider interface {
	// GetID 返回验证器的唯一标识符，如 "totp"
	GetID() string

	// Verify 负责执行验证逻辑
	Verify(ctx context.Context, secret string, code string) (bool, error)
}
