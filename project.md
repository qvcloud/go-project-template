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
*   **数据库迁移**: [sql-migrate](https://github.com/rubenv/sql-migrate) - 功能强大的 SQL 迁移工具，支持单独的配置文件。
*   **可观测性**: [Prometheus](https://prometheus.io/) (指标) & [OpenTelemetry](https://opentelemetry.io/) (追踪)。

## 2. 目录结构 (Directory Structure)

```text
.
├── cmd/                # 应用入口 (Cobra 命令定义)
│   ├── root.go         # 根命令与全局配置初始化
│   └── http.go         # HTTP 服务启动命令
├── config/             # 配置文件与静态资源
│   └── config.yaml     # 运行时配置
├── internal/           # 业务逻辑 (外部不可调用)
│   ├── di/             # 依赖注入中心
│   │   ├── provider/   # 基础架构 Provider (DB, Redis, Log等)
│   │   ├── root.go     # 公共依赖聚合
│   │   └── http.go     # HTTP 服务的 FX 容器定义
│   ├── migrations/     # SQL 迁移文件 (Schema 定义)
│   ├── middleware/     # Gin 中间件 (JWT, 日志平滑, 限流)
│   ├── model/          # 数据模型 (GORM 实体)
│   ├── repository/     # 持久层 (数据库操作)
│   ├── service/        # 业务逻辑层核心
│   └── app/            # 应用入口 (HTTP Handler, MQ Consumer 等各协议实现)
├── pkg/                # 公共工具库 (可被外部引用)
│   ├── response/       # 统一响应格式与错误码定义
│   └── utils/          # 通用辅助函数
├── scripts/            # 运维与开发脚本 (编译、镜像构建等)
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
    - 配置日志模式与连接限制。

### 3.3 日志 (Zap)
*   **Provider**: `NewZapLogger`
*   **职责**: 提供结构化日志，支持 `Console` 和 `File` (带切割) 两种输出。

### 3.4 Web 服务 (Gin)
*   **Provider**: `NewGin`
*   **职责**: 初始化 Gin 引擎，设置中间件（CORS, Recover, Logger）。

### 3.5 声明式路由与中间件 (Router & Middleware)
*   **Router**: `NewRouter` 负责将各 Handler 挂载到 Gin 路由组中。
*   **Middleware**: 包含 JWT 校验、CORS、日志追踪、Recover 等标准中间件。

### 3.6 数据库迁移 (Migration)
*   **工具**: 集成 `sql-migrate`。
*   **职责**: 使用 `.sql` 文件定义迁移，支持 `up` 和 `down` 操作。配置文件通常为 `dbconfig.yml`。

### 3.7 异步消息 (RocketMQ)
*   **Provider**: `NewRocketMQProducer`
*   **职责**: 初始化生产者，在 `OnStop` 时优雅关闭连接。

### 3.8 生命周期 Hook (FX Lifecycle)
框架利用 `fx.Hook` 在 `OnStart` 时启动 HTTP 服务监听，在 `OnStop` 时平滑关闭数据库连接、消息队列与服务器。

## 4. 典型开发流程 (Standard Workflow)

1.  **定义模型**: 在 `internal/model/` 下创建新表结构。
2.  **持久层 (Repository)**:
    - 在 `internal/repository/` 下定义接口并实现。
    - 在 `internal/repository/module.go` 中通过 `fx.Provide` 注册。
3.  **业务层 (Service)**:
    - 在 `internal/service/` 下编写业务逻辑，并注入 Repository。
4.  **控制层 (Handler)**:
    - 在 `internal/app/http/handler/` 下编写处理函数，并注入 Service。
    - 在 `internal/app/http/module.go` 中注册 Provider。
5.  **路由注册**: 在 `internal/app/http/routes.go` 中通过 `Router.Register()` 将路由注册到 `gin.Engine`。

## 5. 项目启动与命令行 (CLI Entrypoints)

通过 `Cobra` 提供的入口启动 HTTP 服务：

*   **启动服务**: `go run main.go http`

## 6. 核心骨架实现参考 (Core Implementation Reference)

### 6.1 应用入口 (Root & CLI)

```go
// cmd/root.go
func initViper() {
    godotenv.Load()
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        viper.AddConfigPath("./config")
        viper.SetConfigName("config")
        viper.SetConfigType("yaml")
    }
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    viper.AutomaticEnv()
    viper.ReadInConfig()
}
```

### 6.2 依赖注入 - 基础 Provider (DI Providers)

#### Viper & Config
```go
// internal/di/provider/viper.provider.go
func NewViper() *viper.Viper {
    return viper.GetViper()
}

// internal/di/provider/config.provider.go
type Env string

const (
    EnvDev  Env = "dev"
    EnvTest Env = "test"
    EnvProd Env = "prod"
)

type Config struct {
    Env      Env            `mapstructure:"env"`
    Log      LogConfig      `mapstructure:"log"`
    HTTP     HTTPConfig     `mapstructure:"http"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    RocketMQ RocketMQConfig `mapstructure:"rocketmq"`
}

type RocketMQConfig struct {
    Endpoint string `mapstructure:"endpoint"`
}

type LogConfig struct {
    Level      string `mapstructure:"level"`
    Format     string `mapstructure:"format"`
    Output     string `mapstructure:"output"`
    File       string `mapstructure:"file"`
    MaxSize    int    `mapstructure:"max_size"`
    MaxBackups int    `mapstructure:"max_backups"`
    MaxAge     int    `mapstructure:"max_age"`
    Compress   bool   `mapstructure:"compress"`
}

type HTTPConfig struct {
    Host         string `mapstructure:"host"`
    Port         int    `mapstructure:"port"`
    ReadTimeout  int    `mapstructure:"read_timeout"`
    WriteTimeout int    `mapstructure:"write_timeout"`
    Mode         string `mapstructure:"mode"`
    JWTSecret    string `mapstructure:"jwt_secret"`
}

type RedisConfig struct {
    URL      string `mapstructure:"url"`
    PoolSize int    `mapstructure:"pool_size"`
    TLS      string `mapstructure:"tls"`
}

type DatabaseConfig struct {
    Driver          string `mapstructure:"driver"` // 固定使用 postgres
    Host            string `mapstructure:"host"`
    Port            int    `mapstructure:"port"`
    Username        string `mapstructure:"username"`
    Password        string `mapstructure:"password"`
    Database        string `mapstructure:"database"`
    MaxOpenConns    int    `mapstructure:"max_open_conns"`
    MaxIdleConns    int    `mapstructure:"max_idle_conns"`
    ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
    Debug           bool   `mapstructure:"debug"`
    TLS             string `mapstructure:"tls"` // 用于 sslmode
}

func NewConfig(v *viper.Viper) (*Config, error) {
    var config Config
    if err := v.Unmarshal(&config); err != nil {
        return nil, err
    }
    return &config, nil
}
```

#### Zap Logger (日志)
```go
// internal/di/provider/zap.provider.go
func NewZapLogger(config *Config) (*zap.Logger, error) {
    cfg := zap.NewProductionConfig()
    cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    if level, err := zapcore.ParseLevel(config.Log.Level); err == nil {
        cfg.Level = zap.NewAtomicLevelAt(level)
    }
    var core zapcore.Core
    var encoder zapcore.Encoder
    if config.Log.Format == "text" {
        cfg.Encoding = "console"
        cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
        encoder = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
    } else {
        encoder = zapcore.NewJSONEncoder(cfg.EncoderConfig)
    }
    if config.Log.Output == "file" && config.Log.File != "" {
        os.MkdirAll(filepath.Dir(config.Log.File), 0755)
        w := zapcore.AddSync(&lumberjack.Logger{
            Filename: config.Log.File, MaxSize: config.Log.MaxSize,
            MaxBackups: config.Log.MaxBackups, MaxAge: config.Log.MaxAge,
            Compress: config.Log.Compress,
        })
        core = zapcore.NewCore(encoder, w, cfg.Level)
    } else {
        core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), cfg.Level)
    }
    return zap.New(core, zap.AddCaller()), nil
}
```

#### GORM (数据库 - Postgres)
```go
// internal/di/provider/gorm.provider.go
func NewGorm(cfg *Config, logger *zap.Logger) (*gorm.DB, error) {
    sslMode := "disable"
    if cfg.Database.TLS != "" {
        sslMode = cfg.Database.TLS
    }
    // 构造 Postgres DSN
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
        cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, 
        cfg.Database.Password, cfg.Database.Database, sslMode)
        
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect database: %w", err)
    }
    
    sqlDB, _ := db.DB()
    sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
    sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
    sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)
    
    if cfg.Database.Debug {
        db = db.Debug()
    }
    return db, nil
}
```

#### Gin (Web 引擎)
```go
// internal/di/provider/gin.provider.go
func NewGin(cfg *Config) *gin.Engine {
    if strings.ToLower(cfg.HTTP.Mode) == "release" {
        gin.SetMode(gin.ReleaseMode)
    }
    return gin.New()
}
```

### 6.3 依赖注入 - 模块聚合 (DI Modules)

#### 公共依赖聚合
```go
// internal/di/root.go
var common = fx.Option(
    fx.Provide(
        provider.NewViper,
        provider.NewConfig,
        provider.NewZapLogger,
        provider.NewGorm,
    ),
)
```

#### HTTP 子应用入口
```go
// internal/di/http.go
func HTTPModule(_ *cobra.Command, _ []string) *fx.App {
    return fx.New(
        common,
        repository.Module,
        service.Module,
        app_http.Module, // 对应 internal/app/http
        fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger) {
            lc.Append(fx.Hook{
                OnStart: func(_ context.Context) error {
                    logger.Info("HTTP Service started")
                    return nil
                },
                OnStop: func(_ context.Context) error {
                    logger.Info("HTTP Service stopped")
                    return nil
                },
            })
        }),
    )
}
```

### 6.4 HTTP 服务生命周期管理

```go
// internal/app/http/module.go (部分)
func setupHTTP(lc fx.Lifecycle, in inFx) {
    server := &http.Server{
        Addr:         fmt.Sprintf(":%d", in.Cfg.HTTP.Port),
        Handler:      in.Engine,
        ReadTimeout:  time.Duration(in.Cfg.HTTP.ReadTimeout) * time.Second,
        WriteTimeout: time.Duration(in.Cfg.HTTP.WriteTimeout) * time.Second,
    }
    lc.Append(fx.Hook{
        OnStart: func(_ context.Context) error {
            go func() {
                in.Router.Register() // 注册路由
                server.ListenAndServe()
            }()
            return nil
        },
        OnStop: func(ctx context.Context) error {
            return server.Shutdown(ctx)
        },
    })
}
```

### 6.5 消息队列 - RocketMQ Producer

```go
// internal/di/provider/rocketmq.provider.go
func NewRocketMQProducer(cfg *Config, lc fx.Lifecycle) (primitive.Producer, error) {
    p, err := rocketmq.NewProducer(
        producer.WithNsResolver(primitive.NewFixedNamesrvResolver([]string{cfg.RocketMQ.Endpoint})),
        producer.WithRetry(2),
    )
    if err != nil {
        return nil, err
    }

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            return p.Start()
        },
        OnStop: func(ctx context.Context) error {
            return p.Shutdown()
        },
    })

    return p, nil
}
```

### 6.6 业务中间件 - JWT 校验与 TraceID

```go
// internal/middleware/trace.go
func TraceID() gin.HandlerFunc {
    return func(c *gin.Context) {
        traceID := c.GetHeader("X-Trace-ID")
        if traceID == "" {
            // 自行实现随机字符串生成
            traceID = "req-" + time.Now().Format("20060102150405")
        }
        c.Set("traceId", traceID)
        c.Header("X-Trace-ID", traceID)
        c.Next()
    }
}

// internal/middleware/auth.go
func JWTAuth(cfg *Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            response.Fail(c, response.CodeUnauthorized, "missing token")
            c.Abort()
            return
        }
        // ... JWT 解析后可将用户信息存入 Context ...
        c.Next()
    }
}
```

### 6.7 请求校验与 Swagger 示例

```go
// internal/app/http/handler/user.go
type LoginReq struct {
    Username string `json:"username" binding:"required,min=4"`
    Password string `json:"password" binding:"required"`
}

// LoginHandler 处理登录请求
// @Summary 用户登录
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body LoginReq true "登录参数"
// @Success 200 {object} response.Response
// @Router /api/v1/login [post]
func LoginHandler(cfg *Config, userService *service.UserService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req LoginReq
        if err := c.ShouldBindJSON(&req); err != nil {
            response.Fail(c, response.CodeInvalidParam, err.Error())
            return
        }
        // ... 调用 service 逻辑 ...
        response.Success(c, nil)
    }
}
```

### 6.8 消费者应用入口 (MQ Consumer)

对于非 HTTP 类型的入口（消息消费者），可以使用独立的 DI 模块启动。

```go
// internal/di/consumer.go
func ConsumerModule(_ *cobra.Command, _ []string) *fx.App {
    return fx.New(
        common,
        repository.Module,
        service.Module,
        fx.Invoke(func(lc fx.Lifecycle, p primitive.Producer, logger *zap.Logger) {
            lc.Append(fx.Hook{
                OnStart: func(_ context.Context) error {
                    logger.Info("MQ Consumer started")
                    // 在此处启动消费组订阅逻辑
                    return nil
                },
                OnStop: func(_ context.Context) error {
                    logger.Info("MQ Consumer stopped")
                    return nil
                },
            })
        }),
    )
}
```

### 6.9 统一响应处理 (Unified Response)
在 `pkg/response/` 中定义标准响应格式，确保前端处理逻辑一致。

```go
// pkg/response/code.go
type Code int

const (
    CodeSuccess      Code = 0
    CodeInvalidParam Code = 400
    CodeUnauthorized Code = 401
    CodeInternalErr  Code = 500
)

// pkg/response/response.go
type Response struct {
    Code    Code        `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{Code: CodeSuccess, Message: "success", Data: data})
}

func Fail(c *gin.Context, code Code, message string) {
    c.JSON(http.StatusOK, Response{Code: code, Message: message})
}
```

## 7. 测试与 Mock 工具 (Mocking & Testing)

项目使用 `go.uber.org/mock` 进行单元测试的 Mock 代码生成，方便在测试 Service 时排除 Repository 或外部 Client 的干扰。

### 7.1 Mock 代码生成
生成的代码统一存放在 `generated/mocks/` 目录下。

*   **生成命令**: `make mock` 或直接运行 `bash scripts/mockgen.sh`。
*   **脚本逻辑**: 内部使用 `mockgen` 根据 `internal/repository` 或其他包中的接口定义自动生成 Mock 实现。

```bash
# scripts/mockgen.sh 示例
mockgen -destination=./generated/mocks/user_repository.go -package=mocks <project_module>/internal/repository UserRepository
```

### 7.2 使用示例
在编写测试时，可以通过 `mocks.NewMockUserRepository(ctrl)` 快速创建一个模拟对象，并预设其返回值。

## 8. 常用脚本工具 (Scripts & Makefile)

项目在 `scripts/` 目录下封装了常用的开发与部署任务，并通过 `Makefile` 提供了统一的命令行入口。

### 8.1 核心脚本说明
*   **`scripts/build.sh`**: 编译脚本。会自动注入版本号、构建时间、Git Commit ID 等元数据到二进制文件中。
*   **`scripts/mockgen.sh`**: 自动化 Mock 生成脚本。
*   **`scripts/build_image.sh`**: Docker 镜像构建脚本。

### 8.2 Makefile 完整命令参考 (Makefile Reference)

为了方便快速搭建环境和执行常用任务，以下是推荐的 `Makefile` 内容，集成了所有必要的工具安装与开发指令：

```makefile
.PHONY: help build run-http clean test docs mock install fmt lint

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## 1. 核心初始化指令：安装项目依赖及所有 CLI 工具
	go mod tidy
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/spf13/cobra-cli@latest
	go install github.com/rubenv/sql-migrate/...@latest

build: ## 编译应用 eg: make build version=v1.0.0
	@bash scripts/build.sh $(version)

run-http: ## 运行 HTTP 服务
	go run main.go http

docs: ## 生成/更新 API 文档 (输出至 generated/docs/)
	swag init -g internal/di/http.go --parseDependency --parseInternal -o ./generated/docs

migrate-up: ## 运行数据库向上迁移
	sql-migrate up -config=config/dbconfig.yml -env="development"

mock: ## 扫描接口定义并刷新 generated/mocks/
	@bash scripts/mockgen.sh

test: ## 运行全项目单元测试 (禁用缓存)
	go test ./... --count=1

fmt: ## 自动代码格式化
	go fmt ./...

lint: ## 严格的高性能静态分析
	golangci-lint run

clean: ## 清理编译产物和 Go 缓存
	rm -rf dist/
	go clean
```

## 9. 核心工具配置文件 (Core Configuration Templates)

以下是项目核心工具的脱敏配置文件，可直接复制到新项目根目录下使用。

### 9.1 代码检查 (`.golangci.yml`)

用于强制执行团队代码规范，集成 `revive`, `gosec`, `staticcheck` 等高性能插件。

```yaml
run:
  timeout: 5m

linters:
  enable:
    - gofmt
    - goimports
    - revive
    - govet
    - staticcheck
    - gosimple
    - ineffassign
    - typecheck
    - unused
    - goconst
    - gosec
    - misspell
    - prealloc
    - errcheck

linters-settings:
  gosec:
    excludes:
      - G104 # 忽略未处理错误检查，由 errcheck 处理
  errcheck:
    check-type-assertions: true
    check-blank: true
    exclude-functions:
      - fmt.Printf
      - fmt.Println
      - github.com/spf13/viper.BindPFlag
      - github.com/joho/godotenv.Load
```

### 9.2 容器化 (`Dockerfile`)

基于 Alpine 的多阶段构建，极致优化镜像体积。

```dockerfile
# 阶段1: 构建
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . /app
RUN go mod download

ARG VERSION
RUN apk add --no-cache make bash git
RUN make build version=${VERSION}

# 阶段2: 运行
FROM alpine:latest
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=build /app/dist/<project_bin> /app/<project_bin>
COPY --from=build /app/config/ /app/config/
RUN chmod +x /app/<project_bin>

EXPOSE 8080
CMD ["./<project_bin>", "http"]
```

### 9.3 基础设施 (`docker-compose.yml`)

快速启动本地开发依赖环境（Postgres, Redis, RocketMQ）。

```yaml
services:
  postgres:
    image: postgres:16-alpine 
    ports:
      - '5432:5432'
    environment:
      - POSTGRES_DB=<project_db>
      - POSTGRES_PASSWORD=postgres

  redis:
    image: redis:latest
    ports:
      - '6379:6379'

  rocketmq-namesrv:
    image: apache/rocketmq:5.3.4
    ports:
      - 9876:9876
    command: sh mqnamesrv

  rocketmq-broker:
    image: apache/rocketmq:5.3.4
    environment:
      - NAMESRV_ADDR=rocketmq-namesrv:9876
    depends_on:
      - rocketmq-namesrv
    ports:
      - 10911:10911
    command: sh mqbroker
```

### 9.4 业务配置 (`config/config.example.yaml`)

包含日志、HTTP、数据库、Redis 等标准配置模板。

```yaml
env: dev
log:
  level: info
  format: json
  output: console # stdout, file, console
  file: logs/app.log

http:
  host: 0.0.0.0
  port: 8080
  mode: debug

database:
  driver: postgres # 固定使用 postgres
  host: localhost
  port: 5432
  username: postgres
  password: postgres
  database: <project_db>
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600
  tls: "disable" # disable, require, verify-ca, verify-full
  debug: false

redis:
  url: "redis://localhost:6379/0"
```

### 9.5 迁移配置 (`config/dbconfig.yml`)

用于 `sql-migrate` 工具的数据库连接配置。

```yaml
development:
  dialect: postgres
  datasource: host=localhost port=5432 user=postgres password=postgres dbname=<project_db> sslmode=disable
  dir: internal/migrations
  table: migrations
```

## 10. 快速复用建议 (Scaffolding Tips)

1.  **一键初始化**: 建议在新项目目录运行以下命令进行重命名：
    ```bash
    find . -type f -not -path '*/.*' -exec sed -i '' 's/go-project-template/<your-project-name>/g' {} +
    ```
2.  **基础设施复用**: 直接复制 `internal/di/provider/` 下的所有文件，它们封装了绝大多数项目的通用依赖（DB, Redis, Log, RocketMQ）。
3.  **DI 骨架**: 复制 `internal/di/root.go` 和 `internal/di/http.go`，以此为基础扩展新命令。
4.  **全局配置**: 修改 `internal/di/provider/config.provider.go` 中的 `Config` 结构体以匹配新业务属性。
5.  **中间件注册**: 在 `internal/app/http/module.go` 的 `NewGin` 中根据需求通过 `r.Use()` 注册中间件。
6.  **统一响应**: 复制 `pkg/response/` 以保持全项目 API 格式的一致性。
