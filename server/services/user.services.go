package services

import (
	"errors"
	"kzchat/server/database"
	"kzchat/server/models"

	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	err := database.DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return user, nil
	}
	return user, err
}

func CheckExistingUser(username string) (bool, error) {
	var user models.User
	result := database.DB.Where("username = ?", username).First(&user)

	if result.Error == nil {
		return true, nil
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, result.Error
}
