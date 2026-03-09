# OpenTelemetry 扩展监控技术设计 (Technical Design)

## 1. 架构演进

我们将从“仅监控基础设施”转向“深度逻辑监控”。核心改动点在于 **Context (上下文)** 的全链路传递。

## 2. 模块设计

### 2.1 Redis 监测 (otelredis)

- **方案**：使用 `github.com/go-redis/redis/extra/redisotel/v8` 官方插件。
- **修改位置**：`app/pkg/data/iredis/client.go`。
- **实现细节**：`
    ```go
    rdb.AddHook(redisotel.NewTracingHook())
    ```

### 2.2 图像处理流水线重构 (重点)

目前 `ImageProcessor` 接口缺乏 `context.Context`，导致追踪中断。

- **修改位置**：`app/pkg/gocv/imgpro/core/process.go`。
- **重构定义**：
    ```go
    type ImageProcessor interface {
        // 增加 ctx 参数
        Process(ctx context.Context, src *gocv.Mat) *gocv.Mat
        GetName() string
    }
    ```
- **埋点实现**：在基类或装饰器中统一处理 Span 开关。
    ```go
    func (p *SomeProcessor) Process(ctx context.Context, src *gocv.Mat) *gocv.Mat {
        ctx, span := otel.Tracer("imgpro").Start(ctx, p.GetName())
        defer span.End()
        // ... 原有逻辑
    }
    ```

### 2.3 Service 层手动埋点规范

- **规范**：每个关键 Service 方法均需使用全局定义的 Tracer。
- **示例**：
    ```go
    func (s *butterflyService) Identify(ctx context.Context, ...) {
        ctx, span := otel.Tracer("service").Start(ctx, "ButterflyService.Identify")
        defer span.End()
        // ...
    }
    ```

## 3. 依赖变更

- 新增：`github.com/go-redis/redis/extra/redisotel/v8`

## 4. 实施路线图

1. **Phase 1**: 集成 `otelredis`，解决缓存监控。
2. **Phase 2**: 重构 `ImageProcessor` 接口，引入 `context.Context`（此步骤涉及文件较多，需谨慎）。
3. **Phase 3**: 在 `app/internal/service` 中全面植入业务 Span。
