package authsvc

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/HunDun0Ben/bs_server/app/internal/service/usersvc"
)

type MFAService struct {
	userSvc *usersvc.UserService
}

func NewMFAService() *MFAService {
	return &MFAService{
		userSvc: usersvc.NewUserService(),
	}
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

func (s *MFAService) VerifyAndActivateMFA(ctx context.Context, username, secret, code string) error {
	// 验证TOTP码
	if !totp.Validate(code, secret) {
		return fmt.Errorf("invalid TOTP code")
	}

	// 生成恢复码
	recoveryCodes, err := s.GenerateRecoveryCodes()
	if err != nil {
		return fmt.Errorf("failed to generate recovery codes: %w", err)
	}

	// 激活MFA并保存恢复码
	return s.userSvc.EnableMFA(ctx, username, secret, recoveryCodes)
}

func (s *MFAService) GetUserMFASecret(ctx context.Context, username string) (string, error) {
	secret, _, err := s.userSvc.GetMFAInfo(ctx, username)
	return secret, err
}

func (s *MFAService) SaveMFASecret(ctx context.Context, username, secret string) error {
	return s.userSvc.SaveMFASecret(ctx, username, secret)
}
