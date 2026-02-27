# 开发任务：FR-01 & FR-02 风险评估与分阶段认证

**相关需求**: FR-01 (风险评估登录), FR-02 (分阶段认证)
**负责人**: [后端开发]

## 1. 任务目标

改造现有登录流程，引入基于风险的认证决策。根据决策结果，系统能颁发两种不同状态的令牌（正式令牌与受限令牌），并通过中间件实现对受限令牌的访问控制。

## 2. 技术实现方案

### 2.1 扩展 JWT Claims

修改 `app/pkg/bsjwt/bs_jwt.go` 中的 `CustomClaims` 结构体，增加 MFA 状态字段。

```go
type CustomClaims struct {
    jwt.RegisteredClaims
    UserID       string `json:"uid"`
    Username     string `json:"unm"`
    // 新增字段
    MFAPending   bool   `json:"mfa_p"`    // 标记是否需要二次验证
    RequiredType string `json:"mfa_type"` // MFA 类型, 如 "totp"
}
```

### 2.2 改造登录逻辑 (`/login`)

在 `app/internal/handler/login_handler.go` 的 `Login` 方法中：

1.  在验证完用户名密码后，**不得**直接生成最终令牌。
2.  **植入风险评估点**：调用一个 `risk_service`（需新建）。
    - **初步实现**：在 `user` 表中增加 `last_login_ip` 字段。该服务比对当前请求 IP 与 `last_login_ip`。若不一致，则判定为高风险。
    - **更新 `last_login_ip`**：无论风险高低，登录成功后都应更新该字段。
3.  **决策与令牌签发**：
    - **低风险**：调用 `bsjwt.GenerateTokenPair`，`MFAPending` 设为 `false`。
    - **高风险**：调用 `bsjwt.GenerateTokenPair`，`MFAPending` 设为 `true`，`RequiredType` 设为 `totp`。并在响应体中增加一个字段 `{"mfa_required": true}`，用于前端判断。

### 2.3 开发 MFA 强制中间件

在 `app/middleware/` 目录下新建 `mfa_enforcer.go`。

- **`MFAEnforcerMiddleware()`**:
    1.  从 `gin.Context` 中获取 JWT `Claims`。
    2.  检查 `Claims.MFAPending` 是否为 `true`。
    3.  若为 `true`，则检查当前请求的 `c.FullPath()` 是否在**白名单**中。
        - **白名单**：应包含 `/api/v1/login/mfa-verify` 和 `/api/v1/logout`。
        - **拦截**：若不在白名单，则调用 `c.AbortWithStatusJSON(403, ...)`，响应体需包含 `{"error": "MFA_REQUIRED", "required_type": "totp"}`。
    4.  若 `MFAPending` 为 `false` 或请求在白名单内，则调用 `c.Next()`。

### 2.4 注册中间件

在 `app/api/router.go` 中，将 `MFAEnforcerMiddleware` 注册到需要 JWT 认证的路由组上，确保它在 `JWTAuth` 中间件之后执行。

```go
auth := apiV1.Group("/")
auth.Use(middleware.JWTAuth())
auth.Use(middleware.MFAEnforcerMiddleware()) // 在 JWT 之后注册
```

### 2.5 新增二次验证接口 (`/login/mfa-verify`)

1.  在 `app/api/router.go` 中为**公开路由组**添加 `POST /login/mfa-verify`。
2.  在 `app/internal/handler/login_handler.go` 中增加 `VerifyMFA` 方法。
    - 该方法暂时留空或只做参数绑定，其具体逻辑将在 `FR-03` 中实现。
    - 需要注意的是，该接口本身**不能**被 `MFAEnforcerMiddleware` 拦截。但它需要一个有效的**受限令牌**才能访问，因此仍需经过 `JWTAuth` 验证。

## 3. 验收标准

- [ ] 高风险登录（如新 IP）后，返回的 JWT 解析后 `MFAPending` 为 `true`。
- [ ] 携带“受限令牌”访问业务接口（如 `/user/insect`）时，返回 403 错误。
- [ ] 携带“受限令牌”访问 `/api/v1/login/mfa-verify` 时，请求应能正常进入 Handler。
- [ ] 低风险登录后，`MFAPending` 为 `false`，所有业务接口可正常访问。
- [ ] 用户模型中 `last_login_ip` 字段被成功更新。
