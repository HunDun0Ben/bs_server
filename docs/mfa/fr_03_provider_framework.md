# 开发任务：FR-03 认证器框架与 TOTP 实现

**相关需求**: FR-03 (认证器框架), F-03 (TOTP 验证支持)
**负责人**: [后端开发]

## 1. 任务目标

设计并实现一个可扩展的 MFA 验证器框架。基于该框架，开发第一个验证器：TOTP（基于时间的一次性密码），并完成二次验证接口的逻辑闭环。

## 2. 技术实现方案

### 2.1 定义 `MFAProvider` 接口

在 `app/internal/service/authsvc/` 目录下新建 `mfa_provider.go`，定义核心接口。

```go
package authsvc

import "context"

// MFAProvider 定义了所有二次验证方式必须实现的接口
type MFAProvider interface {
    // GetID 返回验证器的唯一标识符，如 "totp"
    GetID() string

    // Verify 负责执行验证逻辑
    Verify(ctx context.Context, secret string, code string) (bool, error)
}
```

_注：`secret` 是用户开启 MFA 时存储在数据库中的密钥。_

### 2.2 实现 `TOTPProvider`

在 `app/internal/service/authsvc/` 目录下新建 `totp_provider.go`。

```go
package authsvc

import (
    "context"
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
```

### 2.3 开发 MFA 验证服务

之前的 `MFAService` 过于具体，现在我们创建一个更通用的服务 `MFAVerificationService`，用于管理和调用所有 `Provider`。

在 `app/internal/service/authsvc/` 目录下新建 `mfa_verification_service.go`。

```go
package authsvc

import (
    "context"
    "fmt"
)

// MFAVerificationService 负责管理和调度所有 MFAProvider
type MFAVerificationService struct {
    providers map[string]MFAProvider
}

func NewMFAVerificationService() *MFAVerificationService {
    svc := &MFAVerificationService{
        providers: make(map[string]MFAProvider),
    }
    // 注册所有已实现的 Provider
    svc.register(NewTOTPProvider())
    // svc.register(NewSMSProvider()) // 未来扩展
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
```

### 2.4 完成二次验证接口逻辑 (`/login/mfa-verify`)

回到 `app/internal/handler/login_handler.go` 的 `VerifyMFA` 方法。

1.  **注入依赖**：为 `LoginHandler` 增加 `MFAVerificationService` 和 `UserService` 依赖。
2.  **实现 `VerifyMFA` 逻辑**：
    1.  从 `gin.Context` 获取 `Claims`。
    2.  从请求体中绑定 `otp_code`。
    3.  从 `Claims.RequiredType` 得知本次应使用的 `Provider` 类型。
    4.  调用 `userService.GetMFASecret(claims.UserID)` 获取用户存储的 TOTP 密钥。
    5.  调用 `mfaVerificationService.VerifyCode(providerType, secret, code)` 进行验证。
    6.  **验证成功**:
        - 调用 `bsjwt.GenerateTokenPair` 签发一个**新的、全权限的** JWT（`MFAPending` 设为 `false`）。
        - （可选但推荐）将旧的“受限令牌”的 `jti` 加入 Redis 黑名单，防止重放。
        - 返回新的 `accessToken`。
    7.  **验证失败**：返回 401 Unauthorized。

## 3. 验收标准

- [ ] `MFAProvider` 接口被正确定义。
- [ ] `TOTPProvider` 能够正确验证 `pquerna/otp` 生成的 TOTP 码。
- [ ] `MFAVerificationService` 能够根据类型字符串正确分发到 `TOTPProvider`。
- [ ] 调用 `/login/mfa-verify` 接口：
    - 使用正确的 `otp_code`，能成功返回一个新的、`MFAPending` 为 `false` 的 JWT。
    - 使用错误的 `otp_code`，返回 401 错误。
- [ ] 使用新获取的“正式令牌”可以成功访问所有业务接口。
