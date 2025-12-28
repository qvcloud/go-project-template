package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/qvcloud/go-project-template/internal/delivery/http/handler/user"
	"github.com/qvcloud/go-project-template/internal/di/provider"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	engine *gin.Engine
	logger *zap.Logger
	cfg    *provider.Config
	//middleware
	//controllers
	userHandler *user.Handler
}

type injectContext struct {
	fx.In
	Logger      *zap.Logger
	Engine      *gin.Engine
	UserHandler *user.Handler
	Cfg         *provider.Config
}

func NewHTTPServer(in injectContext) *Server {
	r := &Server{
		logger:      in.Logger,
		cfg:         in.Cfg,
		engine:      in.Engine,
		userHandler: in.UserHandler,
		//各种controller注入进来
	}
	return r
}

func (w *Server) Run() {
	w.logger.Sugar().Info("startup web server")

	w.initRoutes()

	go func() {
		if err := w.engine.Run(fmt.Sprintf("%s:%d", w.cfg.HTTP.Address, w.cfg.HTTP.Port)); err != nil {
			w.logger.Error("failed to start server", zap.Error(err))
		}
	}()
}
