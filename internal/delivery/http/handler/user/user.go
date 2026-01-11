package user

import (
	"github.com/gin-gonic/gin"
	"github.com/qvcloud/go-project-template/internal/service"
	"github.com/qvcloud/go-project-template/pkg/response"
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

// Query godoc
// @Summary Query a user by ID
// @Description Get user details including UUID and nickname
// @Tags User
// @Accept json
// @Produce json
// @Param id query int false "User ID"
// @Success 200 {object} response.Response{data=entity.User}
// @Failure 500 {object} response.Response
// @Router /api/v1/user/query [get]
func (u *Handler) Query(c *gin.Context) {
	// mock id
	var id int64 = 1
	userEntity, err := u.userService.GetUser(id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Success(c, userEntity)
}
