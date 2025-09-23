package service

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"ethereum_task/go/task/task4/common"
	"ethereum_task/go/task/task4/dto"
	"ethereum_task/go/task/task4/repository"
)

func Register(c *gin.Context) {
	var userReq dto.UserReq

	if err := c.ShouldBindJSON(&userReq); err != nil {
		resp := dto.Resp[any]{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
			Data: nil,
		}
		c.JSON(http.StatusOK, resp)
		return
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		resp := dto.Resp[any]{
			Code: http.StatusBadRequest,
			Msg:  "Failed to hash password",
			Data: nil,
		}
		c.JSON(http.StatusOK, resp)
		return
	}
	userReq.Password = string(hashedPassword)

	user := repository.User{
		Email:    userReq.Email,
		Password: userReq.Password,
		UserName: userReq.UserName,
	}

	if err := repository.Register(user); err != nil {
		resp := dto.Resp[any]{
			Code: http.StatusBadRequest,
			Msg:  "Failed to create user",
			Data: nil,
		}
		c.JSON(http.StatusOK, resp)
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, dto.NewRespWithSuccessData(user))
}

func Login(c *gin.Context) {
	var userReq dto.UserReq
	if err := c.ShouldBindJSON(&userReq); err != nil {

		resp := dto.Resp[any]{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
			Data: nil,
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	var storedUser repository.User
	var err error

	if storedUser, err = repository.Login(userReq.UserName, userReq.Password); err != nil {
		resp := dto.Resp[any]{
			Code: http.StatusBadRequest,
			Msg:  "Invalid username or password",
			Data: nil,
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       storedUser.ID,
		"username": storedUser.UserName,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(common.JwtSecret))
	if err != nil {
		resp := dto.Resp[any]{
			Code: http.StatusBadRequest,
			Msg:  "Failed to generate token",
			Data: nil,
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	c.JSON(http.StatusOK, dto.NewRespWithSuccessData("Bearer "+tokenString))

}
