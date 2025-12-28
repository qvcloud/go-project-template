package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/qvcloud/go-project-template/internal/delivery/http/handler/user"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	engine *gin.Engine
	logger *zap.Logger
	viper  *viper.Viper
	//middleware
	//controllers
	userHandler *user.UserHandler
}

type injectContext struct {
	fx.In
	Logger      *zap.Logger
	Viper       *viper.Viper
	Engine      *gin.Engine
	UserHandler *user.UserHandler
}

func NewHTTPServer(in injectContext) *Server {
	r := &Server{
		logger:      in.Logger,
		viper:       in.Viper,
		engine:      in.Engine,
		userHandler: in.UserHandler,
		//各种controller注入进来
	}
	return r
}

func (w *Server) Run() {
	w.logger.Sugar().Info("startup web server")

	w.initRoutes()

	w.viper.SetDefault("listen", "127.0.0.1")
	w.viper.SetDefault("port", 8080)
	go w.engine.Run(fmt.Sprintf("%s:%d", w.viper.GetString("listen"), w.viper.GetInt("port")))
}
