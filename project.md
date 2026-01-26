# 项目框架结构文档 (Project Framework Structure)

本文档详细描述了当前项目的架构设计、目录结构以及核心公共基础 Provider，旨在帮助开发者在新的项目中快速复用和构建相同的框架。

## 1. 核心技术栈 (Core Tech Stack)

*   **编程语言**: Go 1.24+
*   **依赖注入 (DI)**: [Uber-FX](https://github.com/uber-go/fx) - 核心框架，负责所有组件的生命周期管理和依赖分析。
*   **Web 框架**: [Gin](https://github.com/gin-gonic/gin) - 处理 HTTP 请求及路由。
*   **数据库 ORM**: [GORM](https://gorm.io/) - 支持 MySQL, PostgreSQL, SQLite 等。
*   **配置管理**: [Viper](https://github.com/spf13/viper) & [godotenv](https://github.com/joho/godotenv) - 支持多环境配置与环境变量覆盖。
*   **命令行工具**: [Cobra](https://github.com/spf13/cobra) - 构建多子命令应用（如 HTTP 服务、扫链服务、迁移工具）。
*   **日志系统**: [Zap](https://github.com/uber-go/zap) & [Lumberjack](https://github.com/natefinch/lumberjack) - 高性能结构化日志及文件切割。
*   **消息中间件**: [RocketMQ](https://rocketmq.apache.org/) - 异步事件处理与解耦。
*   **API 文档**: [Swagger (swag)](https://github.com/swaggo/swag) - 自动生成 API 接口文档。

## 2. 目录结构 (Directory Structure)

```text
.
├── cmd/                # 应用入口 (Cobra 命令定义)
│   ├── root.go         # 根命令与全局配置初始化
│   ├── http.go         # API 服务启动命令
│   ├── scanner.go      # 区块链解析服务启动命令
│   └── migrate.go      # 数据库迁移工具
├── config/             # 配置文件与静态资源
│   ├── config.yaml     # 运行时配置
│   └── abi/            # 智能合约 ABI JSON 文件
├── internal/           # 业务逻辑 (外部不可调用)
│   ├── di/             # 依赖注入中心
│   │   ├── provider/   # 基础架构 Provider (DB, Redis, Log等)
│   │   ├── root.go     # 公共依赖聚合
│   │   └── http.go     # HTTP 服务的 FX 容器定义
│   ├── model/          # 数据模型 (GORM 实体)
│   ├── repository/     # 持久层 (数据库操作)
│   ├── service/        # 业务逻辑层核心
│   ├── server/         # 接口定义 (HTTP, WebSocket 等)
│   └── migrations/     # 数据库迁移脚本
├── pkg/                # 公共工具库 (可被外部引用)
│   ├── abi_helper/     # 合约解析工具
│   ├── response/       # 统一响应格式定义
│   ├── utils/          # 通用辅助函数
│   └── broker-rocketmq/# 消息队列封装
├── generated/          # 自动生成代码 (Swagger, Mocks)
├── Makefile            # 编译与常用工具指令
└── go.mod              # 依赖管理
```

## 3. 核心基础 Provider (Core Providers)

项目采用 `FX` 提供的依赖注入模式，核心 Provider 定义在 `internal/di/provider/` 中。

### 3.1 配置 (Viper & Config)
*   **Viper**: 初始化配置引擎，支持从 `config.yaml` 或 `.env` 读取。
*   **Config**: 强类型配置结构体，通过 `viper.Unmarshal` 填充，并作为全局依赖注入到各层。

### 3.2 数据库 (GORM)
*   **Provider**: `NewGorm`
*   **职责**:
    - 初始化数据库连接池。
    - 自动检测并执行数据库迁移 (`migrations.Run`)。
    - 配置日志模式与连接限制。

### 3.3 日志 (Zap)
*   **Provider**: `NewZapLogger`
*   **职责**: 提供结构化日志，支持 `Console` 和 `File` (带切割) 两种输出。

### 3.4 Web 服务 (Gin)
*   **Provider**: `NewGin`
*   **职责**: 初始化 Gin 引擎，设置中间件（CORS, Recover, Logger）。

### 3.5 消息队列 (Broker/RocketMQ)
*   **Provider**: `NewBroker`
*   **职责**: 基于 `RocketMQ` 的生产者与消费者封装，用于跨服务或模块间的异步通信。

### 3.6 声明式路由与中间件 (Router & Middleware)
*   **Router**: `NewRouter` 负责将各 Handler 挂载到 Gin 路由组中。
*   **Middleware**: 包含 JWT 校验、CORS、日志追踪等标准中间件。

### 3.7 生命周期 Hook (FX Lifecycle)
框架利用 `fx.Hook` 在 `OnStart` 时启动服务（如 `http.ListenAndServe`、消息订阅、任务处理器等），在 `OnStop` 时平滑关闭连接（如 DB Close、Redis Close）。

## 4. 典型开发流程 (Standard Workflow)

1.  **定义模型**: 在 `internal/model/` 下创建新表结构。
2.  **创建迁移**: 在 `internal/migrations/` 下添加迁移 SQL 或逻辑。
3.  **持久层 (Repository)**:
    - 在 `internal/repository/` 下定义接口并实现。
    - 在 `internal/repository/module.go` 中通过 `fx.Provide` 注册。
4.  **业务层 (Service)**:
    - 在 `internal/service/` 下编写业务逻辑，并注入 Repository。
    - 按功能拆分 Module（如 `event.Module`, `http.Module`）。
5.  **控制层 (Handler)**:
    - 在 `internal/server/http/handler/` 下编写处理函数，并注入 Service。
    - 在 `internal/server/http/module.go` 中注册 Provider。
6.  **路由注册**: 在 `internal/server/http/routes.go` 中通过 `Router.Register()` 将路由注册到 `gin.Engine`。

## 5. 项目启动与命令行 (CLI Entrypoints)

通过 `Cobra` 提供的入口，项目可以以不同的模式运行：

*   **API 服务**: `go run main.go http`
*   **扫链服务**: `go run main.go scanner`
*   **数据库迁移**: `go run main.go migrate`

## 6. 快速复用建议 (Scaffolding Tips)

1.  **基础设施复用**: 直接复制 `internal/di/provider/` 下的所有文件，它们封装了绝大多数项目的通用依赖（DB, Redis, Log）。
2.  **DI 骨架**: 复制 `internal/di/root.go` 和 `internal/di/http.go`，以此为基础扩展新命令。
3.  **全局配置**: 修改 `internal/di/provider/config.provider.go` 中的 `Config` 结构体以匹配新业务属性。
4.  **依赖注入模式**: 始终坚持 `Module` -> `Provide` -> `Invoke` (setup) 的 fx 范式。
5.  **统一响应**: 复制 `pkg/response/` 以保持全项目 API 格式的一致性。
