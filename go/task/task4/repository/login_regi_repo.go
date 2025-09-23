package repository

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Register(user User) error {
	if err := GetDB().Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func Login(userName string, pwd string) (User, error) {
	var storedUser User

	if err := GetDB().Where("user_name = ?", userName).First(&storedUser).Error; err != nil {
		return storedUser, fmt.Errorf("failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(pwd)); err != nil {
		return storedUser, fmt.Errorf("invalid username or password")
	}

	return storedUser, nil
}
