package provider

import (
	"github.com/gin-gonic/gin"
	"github.com/qvcloud/go-project-template/internal/delivery/http/middleware"
)

func NewGin() *gin.Engine {
	engine := gin.Default()
	engine.Use(middleware.CORS())
	return engine
}
