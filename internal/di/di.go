package di

import (
	"context"

	"github.com/qvcloud/go-project-template/internal/delivery/http"
	"github.com/qvcloud/go-project-template/internal/di/provider"
	"github.com/qvcloud/go-project-template/internal/persistence"
	"github.com/qvcloud/go-project-template/internal/usecase"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func App() *fx.App {
	return fx.New(
		fx.Provide(
			zap.NewDevelopment,
			provider.NewViper,
			provider.NewRedis,
			provider.NewGin,
			usecase.NewUserUseCase,
		),
		persistence.Module,
		http.Module,
		fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					logger.Info("application is starting...")
					return nil
				},
				OnStop: func(context.Context) error {
					logger.Info("application is stopping...")
					return nil
				},
			})
		}),
	)
}
