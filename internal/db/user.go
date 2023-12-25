package db

import (
	"todo-list-service/internal/models"
)

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	result := DB.Where("username = ?", username).First(&user)
	return user, result.Error
}
