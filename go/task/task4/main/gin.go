package main

import (
	"github.com/gin-gonic/gin"

	"ethereum_task/go/task/task4/service"
)

func main() {

	router := gin.Default()
	router.Use(service.GlobalErrorHandler())

	//注册
	router.POST("/register", service.Register)
	//登录
	router.POST("/login", service.Login)

	//登录后才能访问
	r2 := router.Group("/api")
	r2.Use(service.JWTAuthMiddleware())
	{
		//测试
		r2.GET("/xx", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message": "Hello World!",
			})

		})

	}

	router.Run() // 默认监听 0.0.0.0:8080

}
