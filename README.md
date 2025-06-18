├── app                               # 应用主目录
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
├── dockerfile                        # Docker 镜像构建脚本
├── go.mod                            
└── go.sum                            
