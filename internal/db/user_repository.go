package db

import (
	"gorm.io/gorm"
	"time"
	"todo-list-service/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetAllUsers(limit int, offset int) ([]models.User, error)
	GetUsersCount() (int64, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateLastLogin(user *models.User) error
}

type GormUserRepository struct {
	DB *gorm.DB
}

func (r *GormUserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.DB.Where("username = ?", username).First(&user)
	return &user, result.Error
}

func (r *GormUserRepository) UpdateLastLogin(user *models.User) error {
	result := r.DB.Model(&user).Update("last_login_at", time.Now())
	return result.Error
}

func (r *GormUserRepository) CreateUser(user *models.User) (*models.User, error) {
	result := r.DB.Create(user)
	return user, result.Error
}

func (r *GormUserRepository) GetAllUsers(limit int, offset int) ([]models.User, error) {
	var users []models.User
	result := r.DB.Limit(limit).Offset(offset).Find(&users)
	return users, result.Error
}

func (r *GormUserRepository) GetUsersCount() (int64, error) {
	var count int64
	result := r.DB.Model(&models.User{}).Count(&count)
	return count, result.Error
}
