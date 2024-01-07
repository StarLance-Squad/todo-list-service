package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"todo-list-service/internal/models"
)

type TodoRepository interface {
	CreateTodo(todo *models.Todo) (*models.Todo, error)
	GetAllTodosForUser(userId float64, limit int, offset int) ([]TodoResponse, error)
	GetTodosCount(float64) (int64, error)
	DeleteTodoByIDAndUserID(todoID string, userID uint) error
	UpdateTodoByIDAndUserID(todoID string, userID uint, updatedTodo *models.Todo) error
}

type GormTodoRepository struct {
	DB *gorm.DB
}

type TodoResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	UserID      float64   `json:"userId"`
}

func (r *GormTodoRepository) CreateTodo(user *models.Todo) (*models.Todo, error) {
	result := r.DB.Create(user)
	return user, result.Error
}

func (r *GormTodoRepository) GetTodosCount(userID float64) (int64, error) {
	var count int64
	result := r.DB.Where("user_id = ?", userID).Model(&models.Todo{}).Count(&count)
	return count, result.Error
}

func (r *GormTodoRepository) GetAllTodosForUser(userId float64, limit int, offset int) ([]TodoResponse, error) {
	var todos []TodoResponse
	result := r.DB.Model(&models.Todo{}).Select("id, title, description, completed, user_id, created_at, updated_at").Where("user_id = ?", userId).Limit(limit).Offset(offset).Scan(&todos)
	return todos, result.Error
}

func (r *GormTodoRepository) DeleteTodoByIDAndUserID(todoID string, userID uint) error {
	result := r.DB.Where("id = ? AND user_id = ?", todoID, userID).Delete(&models.Todo{})
	return result.Error
}

func (r *GormTodoRepository) UpdateTodoByIDAndUserID(todoID string, userID uint, updatedTodo *models.Todo) error {
	result := r.DB.Where("id = ? AND user_id = ?", todoID, userID).Updates(updatedTodo)
	return result.Error
}
