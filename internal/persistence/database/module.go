package database

import (
	"github.com/qvcloud/go-project-template/internal/persistence/database/entities"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Module = fx.Module(
	"database",
	fx.Provide(
		NewUserRepository,
		//注册其他的数据库实例
	),

	fx.Invoke(autoMigrate),
)

func autoMigrate(db *gorm.DB, logger *zap.Logger) error {
	err := db.AutoMigrate(
		&entities.User{},
	)
	if err != nil {
		logger.Error("auto migrate error", zap.Error(err))
		return err
	}
	return nil
}
