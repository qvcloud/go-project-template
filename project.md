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
│   └── http.go         # HTTP 服务启动命令
├── config/             # 配置文件与静态资源
│   └── config.yaml     # 运行时配置
├── internal/           # 业务逻辑 (外部不可调用)
│   ├── di/             # 依赖注入中心
│   │   ├── provider/   # 基础架构 Provider (DB, Redis, Log等)
│   │   ├── root.go     # 公共依赖聚合
│   │   └── http.go     # HTTP 服务的 FX 容器定义
│   ├── model/          # 数据模型 (GORM 实体)
│   ├── repository/     # 持久层 (数据库操作)
│   ├── service/        # 业务逻辑层核心
│   └── server/         # 接口定义 (HTTP 处理函数与路由)
├── pkg/                # 公共工具库 (可被外部引用)
│   ├── response/       # 统一响应格式定义
│   └── utils/          # 通用辅助函数
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
*   **Middleware**: 包含 JWT 校验、CORS、日志追踪等标准中间件。

### 3.6 生命周期 Hook (FX Lifecycle)
框架利用 `fx.Hook` 在 `OnStart` 时启动 HTTP 服务监听，在 `OnStop` 时平滑关闭数据库连接与服务器。

## 4. 典型开发流程 (Standard Workflow)

1.  **定义模型**: 在 `internal/model/` 下创建新表结构。
2.  **持久层 (Repository)**:
    - 在 `internal/repository/` 下定义接口并实现。
    - 在 `internal/repository/module.go` 中通过 `fx.Provide` 注册。
3.  **业务层 (Service)**:
    - 在 `internal/service/` 下编写业务逻辑，并注入 Repository。
4.  **控制层 (Handler)**:
    - 在 `internal/server/http/handler/` 下编写处理函数，并注入 Service。
    - 在 `internal/server/http/module.go` 中注册 Provider。
5.  **路由注册**: 在 `internal/server/http/routes.go` 中通过 `Router.Register()` 将路由注册到 `gin.Engine`。

## 5. 项目启动与命令行 (CLI Entrypoints)

通过 `Cobra` 提供的入口启动 HTTP 服务：

*   **启动服务**: `go run main.go http`

## 7. 核心骨架实现参考 (Core Implementation Reference)

### 7.1 应用入口 (Root & CLI)

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

### 7.2 依赖注入 - 基础 Provider (DI Providers)

#### Viper & Config
```go
// internal/di/provider/viper.provider.go
func NewViper() *viper.Viper {
    return viper.GetViper()
}

// internal/di/provider/config.provider.go
func NewConfig(v *viper.Viper) (*Config, error) {
    var config Config
    setDefaults(v)
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

#### GORM (数据库)
```go
// internal/di/provider/gorm.provider.go
func NewGorm(cfg *Config, logger *zap.Logger) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, 
        cfg.Database.Port, cfg.Database.Database)
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    sqlDB, _ := db.DB()
    sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
    sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
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

### 7.3 依赖注入 - 模块聚合 (DI Modules)

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
        http.Module,
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

### 7.4 HTTP 服务生命周期管理

```go
// internal/server/http/module.go (部分)
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

### 7.5 统一响应处理 (Unified Response)
在 `pkg/response/` 中定义标准响应格式，确保前端处理逻辑一致。

```go
// pkg/response/response.go
type Response struct {
    Code    Code        `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: data})
}

func Fail(c *gin.Context, code Code, message string) {
    c.JSON(http.StatusOK, Response{Code: code, Message: message})
}
```

## 8. 快速复用建议 (Scaffolding Tips)

1.  **基础设施复用**: 直接复制 `internal/di/provider/` 下的所有文件，它们封装了绝大多数项目的通用依赖（DB, Redis, Log）。
2.  **DI 骨架**: 复制 `internal/di/root.go` 和 `internal/di/http.go`，以此为基础扩展新命令。
3.  **全局配置**: 修改 `internal/di/provider/config.provider.go` 中的 `Config` 结构体以匹配新业务属性。
4.  **依赖注入模式**: 始终坚持 `Module` -> `Provide` -> `Invoke` (setup) 的 fx 范式。
5.  **统一响应**: 复制 `pkg/response/` 以保持全项目 API 格式的一致性。
