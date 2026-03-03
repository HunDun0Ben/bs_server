# slog 与 OpenTelemetry 链路追踪集成需求文档 (PRD)

## 1. 项目背景

`bs_server` 作为一个涉及复杂图像处理流水线（Pipeline）的系统，在生产环境中定位问题（如识别失败、耗时过长）具有挑战性。为了实现全链路可观测性，我们需要将**结构化日志（Logging）**与**分布式链路追踪（Tracing）**深度关联。

## 2. 业务目标

- **实现日志与链路关联**：每一行业务日志必须包含当前的 `trace_id` 和 `span_id`。
- **提升运维效率**：开发者可以通过 Jaeger 或 Grafana 看到一条 Trace 的同时，点击查看该请求产生的所有相关日志。
- **标准化日志输出**：统一使用 Go 官方的 `log/slog` 并输出 JSON 格式，方便 ELK 或 Loki 采集。

## 3. 功能需求

- **自适应 Trace 提取**：日志组件应能自动从 `context.Context` 中提取 OTel 注入的 Trace 信息。
- **全局生效**：系统初始化后，所有调用 `slog` 的标准方法都应自动具备 Trace 关联能力。
- **非侵入性**：底层 Service 和 Repository 层不需要感知 Trace 提取的逻辑，只需正确传递 Context。

## 4. 非功能需求

- **性能**：Trace 信息的提取不应导致明显的响应延迟。
- **兼容性**：必须支持 Gin 框架及 MongoDB 驱动产生的 Trace。
