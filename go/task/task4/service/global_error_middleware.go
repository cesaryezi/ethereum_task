package service

import (
	"ethereum_task/go/task/task4/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				log.Printf("Panic recovered: %v", err)

				// 返回统一错误响应
				resp := dto.Resp[any]{
					Code: http.StatusInternalServerError,
					Msg:  "Internal server error",
					Data: nil,
				}
				c.JSON(http.StatusInternalServerError, resp)

				// 中止请求处理
				c.Abort()
			}
		}()

		// 继续处理请求
		c.Next()
	}
}
