package db

import (
	"gorm.io/gorm"
	"time"
	"todo-list-service/internal/models"
)

type TodoRepository interface {
	CreateTodo(todo *models.Todo) (*models.Todo, error)
	GetAllTodosForUser(userId uint, limit int, offset int) ([]TodoResponse, error)
	GetTodosCount() (int64, error)
}

type GormTodoRepository struct {
	DB *gorm.DB
}

type TodoResponse struct {
	ID          string
	Title       string
	Description string
	Completed   bool
	UserID      uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *GormTodoRepository) CreateTodo(user *models.Todo) (*models.Todo, error) {
	result := r.DB.Create(user)
	return user, result.Error
}

func (r *GormTodoRepository) GetTodosCount() (int64, error) {
	var count int64
	result := r.DB.Model(&models.Todo{}).Count(&count)
	return count, result.Error
}

func (r *GormTodoRepository) GetAllTodosForUser(userId uint, limit int, offset int) ([]TodoResponse, error) {
	var todos []TodoResponse
	result := r.DB.Model(&models.Todo{}).Select("id, title, description, completed, user_id, created_at, updated_at").Where("user_id = ?", userId).Limit(limit).Offset(offset).Scan(&todos)
	return todos, result.Error
}
