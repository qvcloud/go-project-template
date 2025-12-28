# Go 项目模板

一个基于整洁架构（Clean Architecture）的现代化 Go Web 服务模板。

## 特性

- **整洁架构 (Clean Architecture)**：实现了 Domain（领域层）、Service（服务层）、Persistence（持久层）和 Delivery（交付层）的关注点分离。
- **依赖注入**: 使用 [Uber Fx](https://github.com/uber-go/fx) 进行依赖管理。
- **Web 框架**: 使用 [Gin](https://github.com/gin-gonic/gin) 构建高性能 HTTP 服务。
- **ORM**: 通过 [GORM](https://gorm.io/) 进行数据库交互（已配置 PostgreSQL 驱动）。
- **配置管理**: 使用 [Viper](https://github.com/spf13/viper) 进行灵活的配置管理。
- **日志**: 使用 [Zap](https://github.com/uber-go/zap) 进行结构化日志记录。
- **缓存**: 集成 [go-redis](https://github.com/redis/go-redis)。
- **API 文档**: 集成 Swagger 自动生成文档。
- **构建工具**: 完善的 Makefile 支持，包括构建、测试、镜像打包等。
- **容器化**: 提供 Dockerfile 支持多阶段构建。

## 项目结构

```text
internal/
├── delivery/           # 交付层 (HTTP, CLI 等外部接口)
│   └── http/           # HTTP 服务与处理器 (Handlers)
├── domain/             # 领域层 (核心业务规则)
│   ├── entity/         # 领域实体
│   └── repository/     # 仓储接口定义
├── service/            # 服务层 (业务逻辑编排)
├── persistence/        # 持久层 (数据库实现细节)
└── di/                 # 依赖注入容器配置
```

## 快速开始

### 前置要求

- Go 1.24+
- PostgreSQL
- Redis

### 安装

1. 克隆仓库：
   ```bash
   git clone https://github.com/qvcloud/go-project-template.git
   cd go-project-template
   ```

2. 安装依赖：
   ```bash
   go mod tidy
   ```

### 配置

应用使用 `viper` 读取配置，默认在根目录下查找 `config.yml` 文件。

`config.yml` 示例：

```yaml
database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  name: "dev"
  debug: true

redis:
  address: "localhost:6379"
  db: 0

listen:
  host: "0.0.0.0"
  port: 8080
```

### 环境变量配置

支持使用环境变量覆盖配置文件中的设置。环境变量前缀为 `APP_`，层级使用 `_` 分隔。

示例：

- `database.host` -> `APP_DATABASE_HOST`
- `redis.address` -> `APP_REDIS_ADDRESS`
- `listen.port` -> `APP_LISTEN_PORT`

### Docker 构建

项目支持构建 Docker 镜像，并允许自定义应用名称和版本：

```bash
# 构建默认镜像 (go-project-template:latest)
make image

# 构建指定名称和版本的镜像
make image APP_NAME=my-app version=v1.0.0
```

### 运行应用

```bash
# 直接运行
go run cmd/main.go

# 使用 Makefile 运行
make run

# 编译二进制文件
make build

# 构建 Docker 镜像
make image

# 运行测试
make test

# 生成 API 文档
make docs

# 代码检查
make lint
```

## 架构说明

本项目遵循 **整洁架构 (Clean Architecture)** 原则：

1.  **领域层 (Domain Layer)** (`internal/domain`): 包含企业业务规则（实体）和应用业务规则（仓储接口）。该层不依赖任何外部库或层。
2.  **服务层 (Service Layer)** (`internal/service`): 负责编排领域实体的数据流向，并指挥实体使用其业务规则来实现用例的目标。
3.  **持久层 (Persistence Layer)** (`internal/persistence`): 实现领域层定义的仓储接口。这是数据库具体实现逻辑（如 GORM 操作）所在的地方。
4.  **交付层 (Delivery Layer)** (`internal/delivery`): 处理外部接口（HTTP, gRPC, CLI）。它负责将外部数据格式转换为用例和实体使用的内部格式。

### 依赖注入

我们使用 `uber-go/fx` 进行依赖注入，DI 容器配置位于 `internal/di`。

- **Providers**: 服务、仓储和基础设施组件的构造函数定义在 `internal/di/provider` 或各自的包中。
- **Modules**: 组件被分组为模块（例如 `persistence.Module`, `http.Module`）。
- **可见性控制**: 数据库连接对象 (`*gorm.DB`) 通过 `fx.Private` 限制在 `persistence` 模块内部使用，强制执行架构边界，防止上层业务逻辑直接依赖数据库实现。

## API 文档

启动服务后，访问 Swagger 文档：

```
http://localhost:8080/swagger/index.html
```

## 许可证

[MIT](LICENSE)
