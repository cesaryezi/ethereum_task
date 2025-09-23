package common

import (
	"net/http"

	"ethereum_task/go/task/task4/dto"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 统一错误处理函数
func ErrorHandler(c *gin.Context, code int, message string) {
	resp := dto.Resp[any]{
		Code: code,
		Msg:  message,
		Data: nil,
	}
	c.JSON(code, resp)
	c.Abort()
}

// BadRequestError 处理400错误
func BadRequestError(c *gin.Context, message string) {
	ErrorHandler(c, http.StatusBadRequest, message)
}

// UnauthorizedError 处理401错误
func UnauthorizedError(c *gin.Context, message string) {
	ErrorHandler(c, http.StatusUnauthorized, message)
}

// ForbiddenError 处理403错误
func ForbiddenError(c *gin.Context, message string) {
	ErrorHandler(c, http.StatusForbidden, message)
}

// NotFoundError 处理404错误
func NotFoundError(c *gin.Context, message string) {
	ErrorHandler(c, http.StatusNotFound, message)
}

// InternalServerError 处理500错误
func InternalServerError(c *gin.Context, message string) {
	ErrorHandler(c, http.StatusInternalServerError, message)
}
