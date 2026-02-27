# `bs_server` 身份安全与自适应 MFA 架构设计文档

## 1. 设计思想 (Design Philosophy)

- **自适应风险控制 (Adaptive Risk Control)**：2FA 不应是全量用户的负担。系统通过风控引擎（如 IP 变更、新设备登录、异常频率）动态评估风险，仅在**中高风险**场景下触发二次验证，实现安全与体验的平衡。
- **认证逻辑与执行策略分离**：Middleware 只负责检查认证的“完备状态”，而不关心认证的具体“实现方式”。
- **状态化令牌降级 (Token Downgrading)**：利用 JWT Claims 传递认证状态。当触发风控时，系统签发一个“受限令牌”（MFAPending），该令牌具有时效性且权限受限。
- **策略模式 (Strategy Pattern)**：支持多验证方式扩展（TOTP, SMS, Email 等），所有验证器遵循统一接口，实现可插拔。

## 2. MFA 业务流程设计 (MFA Process Design)

### 2.1 认证生命周期

1.  **第一阶段（主验证）**：用户通过 `/login` 提交用户名密码。
2.  **风险评估**：后端服务调用风险评估逻辑。
    - **低风险**：正常颁发 `accessToken`，`MFAPending` 设为 `false`。
    - **高风险**：颁发**受限令牌**，`MFAPending` 设为 `true`，并指定 `RequiredType`（如 "totp"）。
3.  **流量拦截**：`MFAEnforcerMiddleware` 识别到 `MFAPending: true`，除白名单外，拦截所有业务请求。
4.  **第二阶段（二次验证）**：前端捕获 `403 MFA_REQUIRED` 状态，引导用户进入 MFA 验证页面。
5.  **验证升级**：调用 `/login/mfa-verify`，验证通过后，服务器撤销受限令牌并颁发全权限令牌。

## 3. 关键模型与接口定义 (Key Definitions)

### 3.1 JWT Claims 状态扩展

在 JWT 负载中增加状态位，用于中间件识别验证进度。

```go
type CustomClaims struct {
    jwt.RegisteredClaims
    UserID       string `json:"uid"`
    // 关键状态位
    MFAPending   bool   `json:"mfa_p"`    // 是否处于等待二次验证状态
    RequiredType string `json:"mfa_type"` // 本次要求的验证方式 (totp/sms/email)
}
```

### 3.2 MFA 验证器抽象 (`MFAProvider`)

通过接口化设计，使系统未来能无缝支持除 TOTP 以外的验证方式。

```go
type MFAProvider interface {
    // GetID 返回验证器标识，如 "totp", "sms"
    GetID() string
    // Verify 执行具体的验证逻辑
    Verify(ctx context.Context, userID string, code string) (bool, error)
    // Send 用于发送类验证（如短信验证码），TOTP 可为空操作
    Send(ctx context.Context, userID string) error
}
```

## 4. 中间件设计 (Middleware Design)

### 4.1 MFA 强制执行中间件 (`MFAEnforcerMiddleware`)

**核心职责**：作为流量关口，阻止未完成二次验证的用户访问敏感资源。

- **拦截逻辑**：
    1.  解析 Context 中的 `Claims`。
    2.  判断 `MFAPending` 是否为 `true`。
    3.  **路径过滤**：
        - 若路径在白名单内（如 `/login/mfa-verify`, `/logout`），则 `c.Next()`。
        - 若不在白名单，则 `c.AbortWithStatusJSON(403, ...)`。
- **响应设计**：返回 `403 Forbidden` 并携带 `required_type` 字段，告知前端应弹出哪种验证框。

## 5. 后续开发指导 (Development Guidance)

### 5.1 统一验证 Handler

建议实现一个统一的 `MFAVerifyHandler`，其内部维护一个注册了所有 `MFAProvider` 的 `map`。当收到请求时，根据 Token 中的 `RequiredType` 自动路由到对应的 Provider。

### 5.2 风险判定逻辑植入

在 `UserService.FindByLogin` 成功后，应植入风控检查点。

- **简单起见**：可先比对当前登录 IP 与数据库中 `last_login_ip` 是否一致。
- **进阶**：引入风险评分系统，根据分值决定是直接放行、触发 MFA、还是直接封禁本次登录。

### 5.3 令牌置换安全

在 `/login/mfa-verify` 验证成功后，务必将原有的受限 Token 加入黑名单或使其失效，并重新颁发具有完整 `Claims` 的新 Token。
