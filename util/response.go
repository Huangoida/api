package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResponseSuccess(c *gin.Context, msg interface{}) {
	resp := map[string]interface{}{
		"code":            200,
		"msg":             msg,
		"requestErrorMsg": nil,
	}
	c.JSON(http.StatusOK, resp)
}

func ResponseError(c *gin.Context, code int, errorType string, msg interface{}) {
	resp := map[string]interface{}{
		"code": code,
		"msg":  nil,
		"requestErrorMsg": map[string]interface{}{
			"errorType": errorType,
			"errorMsg":  msg,
		},
	}
	c.JSON(code, resp)
}
