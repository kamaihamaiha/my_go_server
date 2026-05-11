package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Envelope{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, Envelope{
		Code:    httpStatus,
		Message: message,
		Data:    nil,
	})
}
