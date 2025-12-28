package user

import (
	"github.com/gin-gonic/gin"
	"github.com/qvcloud/go-project-template/internal/usecase"
	"go.uber.org/zap"
)

type UserHandler struct {
	logger      *zap.Logger
	userUseCase usecase.UserUseCase
}

func NewUserHandler(logger *zap.Logger, user usecase.UserUseCase) *UserHandler {
	r := UserHandler{
		logger:      logger,
		userUseCase: user,
	}
	return &r
}

func (u *UserHandler) Query(c *gin.Context) {

}
