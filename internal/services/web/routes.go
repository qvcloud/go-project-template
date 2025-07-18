package web

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (r *WebServ) initRoutes() {
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.engine.GET("/test", func(c *gin.Context) {
		c.String(200, time.Now().String())
	})

	apiV1 := r.engine.Group("/api/v1")
	apiV1.GET("/user/query", r.userController.Query)
}
