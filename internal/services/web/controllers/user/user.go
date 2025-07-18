package user

import (
	"github.com/gin-gonic/gin"
	"github.com/qvcloud/go-project-template/internal/persistence"
	"go.uber.org/zap"
)

type UserController struct {
	logger   *zap.Logger
	userRepo persistence.UserRepository
}

func NewUserController(logger *zap.Logger, user persistence.UserRepository) *UserController {
	r := UserController{
		logger:   logger,
		userRepo: user,
	}
	return &r
}

func (u *UserController) Query(c *gin.Context) {

}
