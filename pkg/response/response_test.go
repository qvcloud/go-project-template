package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/qvcloud/go-project-template/pkg/response"
	"github.com/stretchr/testify/assert"
)

func TestResponseSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.Success(c, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var res response.Response
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.Nil(t, err)
	assert.Equal(t, response.CodeSuccess, res.Code)
	assert.Equal(t, "success", res.Message)
	assert.Nil(t, res.Data)
}

func TestResponseSuccessWithData(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := "hi"
	response.Success(c, data)

	assert.Equal(t, http.StatusOK, w.Code)
	var res response.Response
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.Nil(t, err)
	assert.Equal(t, response.CodeSuccess, res.Code)
	assert.Equal(t, data, res.Data)
}

func TestResponseFail(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.Fail(c, 500, "unknown error")

	assert.Equal(t, http.StatusOK, w.Code)
	var res response.Response
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.Nil(t, err)
	assert.Equal(t, 500, res.Code)
	assert.Equal(t, "unknown error", res.Message)
}

func TestResponseFailWithError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.FailWithError(c, assert.AnError)

	assert.Equal(t, http.StatusOK, w.Code)
	var res response.Response
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.Nil(t, err)
	assert.Equal(t, response.CodeFailUnknown, res.Code)
	assert.Equal(t, assert.AnError.Error(), res.Message)
}
