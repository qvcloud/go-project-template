package provider

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewZapLogger(config *Config) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 设置日志级别
	if level, err := zapcore.ParseLevel(config.Log.Level); err == nil {
		cfg.Level = zap.NewAtomicLevelAt(level)
	}

	var core zapcore.Core

	// 设置日志格式
	var encoder zapcore.Encoder
	if config.Log.Format == "text" {
		cfg.Encoding = "console"
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(cfg.EncoderConfig)
	}

	// 设置输出路径
	if config.Log.Output == "file" && config.Log.File != "" {
		// 确保日志目录存在
		if err := os.MkdirAll(filepath.Dir(config.Log.File), 0755); err != nil {
			return nil, err
		}

		// 使用 lumberjack 进行日志切割
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.Log.File,
			MaxSize:    config.Log.MaxSize,    // megabytes
			MaxBackups: config.Log.MaxBackups, // max backups
			MaxAge:     config.Log.MaxAge,     // days
			Compress:   config.Log.Compress,   // disabled by default
		})

		core = zapcore.NewCore(
			encoder,
			w,
			cfg.Level,
		)
	} else {
		// 默认输出到 stdout
		core = zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			cfg.Level,
		)
	}

	return zap.New(core, zap.AddCaller()), nil
}
