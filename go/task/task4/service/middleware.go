package service

import (
	"ethereum_task/go/task/task4/common"
	"ethereum_task/go/task/task4/dto"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			resp := dto.Resp[any]{
				Code: http.StatusUnauthorized,
				Msg:  "Missing authorization token",
				Data: nil,
			}
			c.JSON(http.StatusOK, resp)
			c.Abort()
			return
		}

		// 如果token以"Bearer "开头，移除前缀
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(common.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			resp := dto.Resp[any]{
				Code: http.StatusUnauthorized,
				Msg:  "Invalid or expired token",
				Data: nil,
			}
			c.JSON(http.StatusOK, resp)
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["id"])
			c.Set("username", claims["username"])
		}
		// 前置处理
		c.Next()
		// 后置处理
	}
}
