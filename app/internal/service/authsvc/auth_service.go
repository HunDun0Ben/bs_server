package authsvc

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/HunDun0Ben/bs_server/app/internal/repository"
	"github.com/HunDun0Ben/bs_server/app/internal/service/usersvc"
)

const (
	BlockedListKeyTemplate  = "blockedList_jti:%s"
	RefreshTokenKeyTemplate = "refresh_token:%s"
)

type AuthService interface {
	StoreRefreshToken(ctx context.Context, jti, username string, expiration time.Duration) error
	IsRefreshTokenValid(ctx context.Context, jti string) (string, error)
	InvalidateAccessToken(ctx context.Context, jti string, expiration time.Duration) error
	InvalidateRefreshToken(ctx context.Context, jti string) error
	IsAccessTokenValid(ctx context.Context, jti string) (bool, error)

	// MFA methods
	GenerateTOTPSecret(username string) (*otp.Key, error)
	GenerateRecoveryCodes() ([]string, error)
	ValidateTOTPCode(secret string, code string) bool
	VerifyAndActivateMFA(ctx context.Context, username, secret, code string) error
	GetUserMFASecret(ctx context.Context, username string) (string, error)
	SaveMFASecret(ctx context.Context, username, secret string) error
}

type authService struct {
	repo    repository.AuthRepository
	userSvc usersvc.UserService
}

func NewAuthService(repo repository.AuthRepository, userSvc usersvc.UserService) AuthService {
	return &authService{
		repo:    repo,
		userSvc: userSvc,
	}
}

func (s *authService) StoreRefreshToken(ctx context.Context, jti, username string, expiration time.Duration) error {
	redisKey := fmt.Sprintf(RefreshTokenKeyTemplate, jti)
	return s.repo.Set(ctx, redisKey, username, expiration)
}

func (s *authService) IsRefreshTokenValid(ctx context.Context, jti string) (string, error) {
	redisKey := fmt.Sprintf(RefreshTokenKeyTemplate, jti)
	return s.repo.Get(ctx, redisKey)
}

func (s *authService) InvalidateAccessToken(ctx context.Context, jti string, expiration time.Duration) error {
	redisKey := fmt.Sprintf(BlockedListKeyTemplate, jti)
	return s.repo.Set(ctx, redisKey, "1", expiration)
}

func (s *authService) InvalidateRefreshToken(ctx context.Context, jti string) error {
	redisKey := fmt.Sprintf(RefreshTokenKeyTemplate, jti)
	size, err := s.repo.Del(ctx, redisKey)
	if err != nil {
		return err
	}
	if size > 0 {
		return nil
	}
	return errors.New("key 不存在")
}

func (s *authService) IsAccessTokenValid(ctx context.Context, jti string) (bool, error) {
	redisKey := fmt.Sprintf(BlockedListKeyTemplate, jti)
	val, err := s.repo.Get(ctx, redisKey)
	if err != nil {
		return false, err
	}
	if val == "" {
		return true, nil
	}
	return false, nil
}

func (s *authService) GenerateTOTPSecret(username string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "HunDun0Ben/bs_server",
		AccountName: username,
	})
}

func (s *authService) GenerateRecoveryCodes() ([]string, error) {
	codes := make([]string, 8)
	for i := 0; i < 8; i++ {
		b := make([]byte, 5)
		if _, err := rand.Read(b); err != nil {
			return nil, err
		}
		codes[i] = strings.ToUpper(base32.StdEncoding.EncodeToString(b))[:8]
	}
	return codes, nil
}

func (s *authService) ValidateTOTPCode(secret string, code string) bool {
	return totp.Validate(code, secret)
}

func (s *authService) VerifyAndActivateMFA(ctx context.Context, username, secret, code string) error {
	if !totp.Validate(code, secret) {
		return fmt.Errorf("invalid TOTP code")
	}
	recoveryCodes, err := s.GenerateRecoveryCodes()
	if err != nil {
		return fmt.Errorf("failed to generate recovery codes: %w", err)
	}
	return s.userSvc.EnableMFA(ctx, username, secret, recoveryCodes)
}

func (s *authService) GetUserMFASecret(ctx context.Context, username string) (string, error) {
	secret, _, err := s.userSvc.GetMFAInfo(ctx, username)
	return secret, err
}

func (s *authService) SaveMFASecret(ctx context.Context, username, secret string) error {
	return s.userSvc.SaveMFASecret(ctx, username, secret)
}
