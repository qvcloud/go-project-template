package persistence

import (
	"fmt"

	"github.com/qvcloud/go-project-template/internal/di/provider"
	"github.com/qvcloud/go-project-template/internal/domain/entity"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Module = fx.Module(
	"persistence",
	fx.Provide(
		provider.NewGorm,
		fx.Private,
	),
	fx.Provide(
		NewUserRepository,
		//注册其他的数据库实例
	),

	fx.Invoke(autoMigrate),
)

func autoMigrate(db *gorm.DB, logger *zap.Logger) error {
	migrator := db.Migrator()

	models := []interface{}{
		&entity.User{},
		//在这里添加其他的实体模型
	}

	for _, model := range models {
		// If table does not exist, create it (safe)
		if !migrator.HasTable(model) {
			if err := migrator.CreateTable(model); err != nil {
				logger.Error("create table failed", zap.Error(err))
				return fmt.Errorf("create table failed: %w", err)
			}
			logger.Info("created table for model", zap.Any("model", fmt.Sprintf("%T", model)))
			continue
		}

		// Table exists: only add missing columns, do not alter or drop existing ones.
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			logger.Error("parse model schema failed", zap.Error(err))
			return fmt.Errorf("parse model schema failed: %w", err)
		}

		for _, field := range stmt.Schema.Fields {
			// field.DBName is the actual column name in DB; AddColumn expects the struct field name
			if !migrator.HasColumn(model, field.DBName) {
				if err := migrator.AddColumn(model, field.Name); err != nil {
					logger.Error("add column failed", zap.String("column", field.DBName), zap.Error(err))
					return fmt.Errorf("add column %s failed: %w", field.DBName, err)
				}
				logger.Info("added missing column", zap.String("column", field.DBName))
			}
		}
	}

	return nil
}
