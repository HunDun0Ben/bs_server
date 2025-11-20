# 蝴蝶识别系统后端 API

这是一个使用 Go、Gin、GoCV 和 MongoDB 构建的蝴蝶识别系统的服务端 API。它支持前端进行模型训练和用户识别。系统提供可配置的管道，允许用户组合不同的特征提取算法、聚类方法和训练策略，以根据配置生成相应的训练器。

## 🚀 核心功能

*   **模型训练:** 支持用户配置不同的算法和策略来训练图像识别模型。
*   **用户识别:** 提供基于已训练模型的图像识别服务。
*   **可配置管道:** 允许灵活组合特征提取、聚类和训练策略。

## 💡 技术栈

*   **后端:** Go (Gin)
*   **图像处理:** GoCV (OpenCV bindings for Go)
*   **数据库:** MongoDB
*   **缓存:** Redis
*   **认证:** JWT (JSON Web Tokens)
*   **配置管理:** Viper
*   **API 文档:** Swagger
*   **分布式追踪:** OpenTelemetry

## 🛠️ 环境要求

*   Go 1.24+
*   MongoDB
*   Redis
*   OpenTelemetry Collector (用于分布式追踪)

## 🚀 快速开始

### 1. 构建项目

使用 `Makefile` 构建主应用程序和所有工具：
```bash
make all
```

### 2. 运行应用程序

构建成功后，执行生成的二进制文件：
```bash
./scripts/bin/bs_server
```
应用程序将根据 `conf/application.yaml` 中配置的端口启动（默认通常为 `8080`）。

## 📚 API 文档

API 文档通过 Swagger 自动生成。

### 生成文档
```bash
make swagger
```

### 查看文档
应用程序运行后，在浏览器中访问 `http://localhost:8080/swagger/index.html` (如果服务器运行在默认端口和 `BasePath`) 即可查看完整的 API 文档。

## 📂 项目结构

```text
├── app
│   ├── api                           # 路由入口层
│   ├── conf                          # 配置及配置模型
│   ├── docs                          # 项目文档
│   ├── internal                      # 内部核心结构
│   │   ├── handler                   # 控制器/接口处理
│   │   ├── model                     # 结构定义（数据/常量）
│   │   ├── repository                # 持久化层（可选）
│   │   └── service                   # 业务服务逻辑
│   ├── main                          # 启动模块
│   ├── middleware                    # 中间件（认证、错误处理等）
│   ├── pkg                           # 通用库与工具包
│   │   ├── conf                      # 配置加载模块
│   │   ├── data                      # 数据访问封装（如 Mongo）
│   │   ├── errors                    # 错误处理
│   │   ├── gocv                      # 图像处理工具
│   │   ├── kafka                     # Kafka 工具封装
│   │   └── util                      # 通用工具集合
│   ├── script                        # 脚本与初始化数据
│   └── test                          # 测试相关代码（预留）
├── dockerfile
├── go.mod
└── go.sum
```