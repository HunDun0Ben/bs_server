# slog 与 OpenTelemetry 链路追踪集成开发文档 (Technical Design)

## 1. 核心架构设计

为了不破坏 `slog` 的原生性能，我们采用了 **装饰器模式 (Decorator Pattern)**。通过实现 `slog.Handler` 接口，在日志写入磁盘/控制台之前，动态地将 Trace 信息注入到 Attribute 中。

### 1.1 关键组件

- **`OTelHandler`**: 位于 `app/pkg/logger`。它拦截 `Handle(ctx, record)` 调用，使用 `go.opentelemetry.io/otel/trace` 包从 `ctx` 提取 Span 状态。
- **Context 传播链路**:
  `otelgin (中间件)` -> `Gin Handler` -> `slog.InfoContext` -> `OTelHandler`

## 2. 代码实现规范

### 2.1 日志组件初始化

在 `app/main.go` 的程序入口处进行全局初始化。必须在 `initTracerProvider` 之前或同时设置，确保启动日志也能被捕获。

```go
// 初始化 JSON 格式并挂载 OTel 拦截器
handler := logger.NewOTelHandler(slog.NewJSONHandler(os.Stdout, nil))
slog.SetDefault(slog.New(handler))
```

### 2.2 上下文传递准则 (重要)

为了保证 Trace ID 不丢失，必须严格遵守以下编码规范：

1.  **禁止在 Service/Repo 层使用 `*gin.Context`**：
    - Service 接口签名应统一使用 `ctx context.Context`。
    - 这保证了业务逻辑与 Web 框架解耦，且符合 RPC/微服务拆分标准。

2.  **在 Handler 层进行“翻译”**：
    - 调用 Service 时，必须传入 `cxt.Request.Context()` 而不是 `cxt` 本身。
    - 原因：`otelgin` 将 Trace 信息注入在 `http.Request` 的上下文对象中。

3.  **使用 Context 感知方法**：
    - 必须使用 `slog.InfoContext(ctx, ...)` 或 `slog.ErrorContext(ctx, ...)`。
    - 普通的 `slog.Info(...)` 会忽略上下文，导致 `trace_id` 丢失。

## 3. 日志结构示例

集成后，输出的 JSON 日志将包含以下字段：

```json
{
    "time": "2026-02-28T14:30:03.123Z",
    "level": "ERROR",
    "msg": "更新用户登录信息失败",
    "error": "connection timeout",
    "trace_id": "229dbb08ccd14147d2eedfabf360bc28",
    "span_id": "d767c4fb0328ca6d"
}
```

## 4. 扩展性说明

- **微服务化**：本方案完全兼容 gRPC。当系统拆分为微服务时，OTel 的 `Propagator` 会自动通过 Header 传递 Context。只需在 RPC Server 接收端同样配置 `OTelHandler`，即可实现跨服务的 Trace 关联。
- **审计日志**：可以通过扩展 `OTelHandler`，自动将 `user_id` 等通用字段也注入日志。

## 5. 验证步骤

1.  启动项目并配置 `application.yaml` 中的 `otel.enable: true`。
2.  调用任一 API（如 `/api/v1/login`）。
3.  检查控制台输出，确认包含 32 位长度的 `trace_id`。
4.  确认该 `trace_id` 与 OpenTelemetry 控制面板中的 Trace ID 一致。
