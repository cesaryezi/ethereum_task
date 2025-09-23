package repository

import (
	"ethereum_task/go/task/task4/logger"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func Register(user User) error {
	if err := GetDB().Create(&user).Error; err != nil {
		logger.Logger.Error("failed to create user", zap.Error(err))
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func Login(userName string, pwd string) (User, error) {
	var storedUser User

	if err := GetDB().Where("user_name = ?", userName).First(&storedUser).Error; err != nil {
		logger.Logger.Error("failed to find user", zap.String("username", userName), zap.Error(err))
		return storedUser, fmt.Errorf("failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(pwd)); err != nil {
		logger.Logger.Warn("invalid username or password", zap.String("username", userName))
		return storedUser, fmt.Errorf("invalid username or password")
	}

	return storedUser, nil
}
