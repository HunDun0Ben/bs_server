package authsvc

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type MFAService struct{}

func NewMFAService() *MFAService {
	return &MFAService{}
}

func (s *MFAService) GenerateTOTPSecret(username string) (*otp.Key, error) {
	// 生成TOTP密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "HunDun0Ben/bs_server",
		AccountName: username,
	})
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (s *MFAService) GenerateRecoveryCodes() ([]string, error) {
	codes := make([]string, 8) // 生成8个恢复码
	for i := 0; i < 8; i++ {
		b := make([]byte, 5)
		if _, err := rand.Read(b); err != nil {
			return nil, err
		}
		codes[i] = strings.ToUpper(base32.StdEncoding.EncodeToString(b))[:8]
	}
	return codes, nil
}

func (s *MFAService) ValidateTOTPCode(secret string, code string) bool {
	return totp.Validate(code, secret)
}

func (s *MFAService) VerifyAndActivateMFA(username, secret, code string) error {
	// 验证TOTP码
	if !totp.Validate(code, secret) {
		return fmt.Errorf("invalid TOTP code")
	}

	// TODO: 在数据库中更新用户的MFA状态为已激活
	// 1. 更新用户的MFA secret
	// 2. 设置MFA状态为已激活
	// 3. 存储恢复码(如果有的话)

	return nil
}

func (s *MFAService) GetUserMFASecret(username string) (string, error) {
	// TODO: 从数据库中获取用户的TOTP secret
	// 如果用户未设置MFA或MFA未激活，返回错误
	return "", nil
}
