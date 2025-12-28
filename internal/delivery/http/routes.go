package http

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) initRoutes() {
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.engine.GET("/test", func(c *gin.Context) {
		c.String(200, time.Now().String())
	})

	apiV1 := s.engine.Group("/api/v1")
	apiV1.GET("/user/query", s.userHandler.Query)
}
