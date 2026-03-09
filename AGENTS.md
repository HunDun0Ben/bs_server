# 项目上下文：bs_server

**由 Sisyphus 生成于 2026-03-01 01:21:38**

## 概览

`bs_server` 是一个用于蝴蝶识别的 Go 后端服务。它使用 Gin 框架处理 API 请求，GoCV 进行图像处理，并采用分层架构以分离业务逻辑和数据持久化。

## 项目结构

```
bs_server/
├── app/              # 所有应用代码的根目录
│   ├── api/          # API 路由定义 (router.go)
│   ├── internal/     # 核心业务逻辑 (Handler, Service, Repo)
│   ├── pkg/          # 可复用的通用代码
│   ├── main.go       # 应用主入口
│   └── scripts/      # 工具脚本
├── conf/             # 配置文件 (application.yaml)
├── Makefile          # 开发任务自动化
└── go.mod            # Go 模块定义
```

## 在哪里寻找

| 任务                | 位置                       | 注意                                |
| ------------------- | -------------------------- | ----------------------------------- |
| **启动服务**        | `app/main.go`              | 初始化配置、数据库并启动 Gin 服务。 |
| **添加新 API 路由** | `app/api/router.go`        | 在 `InitRoute` 函数中注册新路由。   |
| **修改业务逻辑**    | `app/internal/service/`    | 核心业务逻辑所在地。                |
| **处理 HTTP 请求**  | `app/internal/handler/`    | 解析请求、调用 Service 并返回响应。 |
| **数据持久化**      | `app/internal/repository/` | 与数据库（MongoDB, Redis）交互。    |
| **图像处理**        | `app/pkg/gocv/`            | **注意：** 必须手动管理内存。       |
| **修改应用配置**    | `conf/*.yaml`              | 修改服务端口、数据库连接等。        |

## 核心开发约定

- **分层架构**: 严格遵守 `Handler` -> `Service` -> `Repository` 的调用顺序。
- **依赖注入**: 所有依赖通过构造函数注入，禁止全局状态。
- **错误处理**: 统一使用 `bsvo.AppError` 类型进行错误上报，由中间件统一处理。
- **代码格式化**: 在提交前务必运行 `make format`，它会使用 `gofmt` 和 `gci` 统一代码风格和 import 顺序。
- **内存管理**: 在 `app/pkg/gocv` 中处理 `gocv.Mat` 时，必须遵循“谁创建，谁释放 (`defer mat.Close()`)”的原则。
- **Service 链路追踪**: Service 层方法必须将 `context.Context` 作为第一个参数，以确保 Trace 链路不中断。
- **日志关联**: 在 Service 层记录日志时，必须严格使用上下文感知的方法 (如 `slog.InfoContext`)，以保证日志包含 `trace_id`。

## 关键命令

```bash
# 构建所有 (应用 + 工具)
make all

# 仅构建主应用
make build

# 格式化代码 (提交前必做)
make format

# 生成 Swagger 文档
make swagger

# 运行测试
make test
```

## 注意事项

- **非标准布局**: 项目将所有代码放在 `app/` 目录下，这与 Go 社区主流的顶层 `cmd/`, `internal/`, `pkg/` 布局不同。
- **多个 main**: `app/scripts/` 目录下包含多个 `main` 包，它们是独立的工具脚本，而非主服务。
- **本地开发**: `go.mod` 中的 `replace` 指令用于本地开发，CI/CD 环境可能需要调整。
