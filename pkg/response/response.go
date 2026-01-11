package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

const (
	CodeSuccess     = 0
	CodeFailUnknown = 1000
)

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func FailWithError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeFailUnknown,
		Message: err.Error(),
		Data:    nil,
	})
}
