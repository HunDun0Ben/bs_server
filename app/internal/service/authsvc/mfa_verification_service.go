package authsvc

import (
	"context"
	"fmt"
)

// MFAVerificationService 负责管理和调度所有 MFAProvider
type MFAVerificationService struct {
	providers map[string]MFAProvider
}

func NewMFAVerificationService(providers ...MFAProvider) *MFAVerificationService {
	svc := &MFAVerificationService{
		providers: make(map[string]MFAProvider),
	}
	if len(providers) == 0 {
		// 默认注册所有已实现的 Provider
		svc.register(NewTOTPProvider())
	} else {
		for _, p := range providers {
			svc.register(p)
		}
	}
	return svc
}

func (s *MFAVerificationService) register(provider MFAProvider) {
	s.providers[provider.GetID()] = provider
}

// VerifyCode 统一的验证入口
func (s *MFAVerificationService) VerifyCode(ctx context.Context, providerType, secret, code string) (bool, error) {
	provider, ok := s.providers[providerType]
	if !ok {
		return false, fmt.Errorf("unsupported MFA provider type: %s", providerType)
	}
	return provider.Verify(ctx, secret, code)
}
