package provider

import (
	"github.com/gin-gonic/gin"
)

func NewGin() *gin.Engine {
	engine := gin.New()
	return engine
}
