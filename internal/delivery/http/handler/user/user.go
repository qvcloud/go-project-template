package user

import (
	"github.com/gin-gonic/gin"
	"github.com/qvcloud/go-project-template/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	logger      *zap.Logger
	userService service.UserService
}

func NewHandler(logger *zap.Logger, user service.UserService) *Handler {
	r := Handler{
		logger:      logger,
		userService: user,
	}
	return &r
}

func (u *Handler) Query(_ *gin.Context) {

}
