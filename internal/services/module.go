package services

import (
	"context"

	"github.com/qvcloud/go-project-template/internal/services/web"
	"github.com/qvcloud/go-project-template/internal/services/web/controllers/user"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Modules = fx.Module(
	"services",
	fx.Provide(
		//注册各个模块的控制器
		user.NewUserController,
		web.NewWebServ,
	),
	fx.Invoke(run),
)

func run(lc fx.Lifecycle, logger *zap.Logger, v *viper.Viper, web *web.WebServ) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			web.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
