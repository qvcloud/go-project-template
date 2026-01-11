package http

import (
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/qvcloud/go-project-template/generated/docs" // swagger docs
	"github.com/qvcloud/gopkg/version"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
)

func (s *Server) initRoutes() {
	if version.Version != "" {
		if doc := swag.GetSwagger("swagger"); doc != nil {
			if spec, ok := doc.(*swag.Spec); ok {
				spec.Version = version.Version
			}
		}

	}

	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	apiV1 := s.engine.Group("/api/v1")
	apiV1.GET("/user/query", s.userHandler.Query)
}
