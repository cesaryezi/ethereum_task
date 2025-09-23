// service/login_register.go
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
		common.BadRequestError(c, err.Error())
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		common.InternalServerError(c, "Failed to hash password")
		return
	}
	userReq.Password = string(hashedPassword)

	user := repository.User{
		Email:    userReq.Email,
		Password: userReq.Password,
		UserName: userReq.UserName,
	}

	if err := repository.Register(user); err != nil {
		common.BadRequestError(c, "Failed to create user")
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, dto.NewRespWithSuccessData(user))
}

func Login(c *gin.Context) {
	var userReq dto.UserReq
	if err := c.ShouldBindJSON(&userReq); err != nil {
		common.BadRequestError(c, err.Error())
		return
	}

	var storedUser repository.User
	var err error

	if storedUser, err = repository.Login(userReq.UserName, userReq.Password); err != nil {
		common.UnauthorizedError(c, "Invalid username or password")
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
		common.InternalServerError(c, "Failed to generate token")
		return
	}

	c.JSON(http.StatusOK, dto.NewRespWithSuccessData("Bearer "+tokenString))
}
