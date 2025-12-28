package http

import (
	"context"

	"github.com/qvcloud/go-project-template/internal/delivery/http/handler/user"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"http",
	fx.Provide(
		//注册各个模块的控制器
		user.NewUserHandler,
		NewHTTPServer,
	),
	fx.Invoke(run),
)

func run(lc fx.Lifecycle, logger *zap.Logger, v *viper.Viper, s *Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
