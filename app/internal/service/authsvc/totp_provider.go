package authsvc

import (
	"context"
	"errors"

	"github.com/pquerna/otp/totp"
)

// TOTPProvider 实现了 MFAProvider 接口
type TOTPProvider struct{}

func NewTOTPProvider() *TOTPProvider {
	return &TOTPProvider{}
}

func (p *TOTPProvider) GetID() string {
	return "totp"
}

func (p *TOTPProvider) Verify(_ context.Context, secret string, code string) (bool, error) {
	valid := totp.Validate(code, secret)
	if !valid {
		return false, errors.New("invalid TOTP code")
	}
	return true, nil
}
