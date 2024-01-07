package services

import (
	"todo-list-service/internal/db"
	"todo-list-service/internal/models"
)

type TodoService struct {
	TodoRepo db.TodoRepository
}

func (s *TodoService) CreateTodo(todo *models.Todo) (*models.Todo, error) {
	return s.TodoRepo.CreateTodo(todo)
}

func (s *TodoService) GetTodosCount(userId float64) (int64, error) {
	return s.TodoRepo.GetTodosCount(userId)
}

func (s *TodoService) GetAllTodosForUser(userId float64, limit int, offset int) ([]db.TodoResponse, error) {
	return s.TodoRepo.GetAllTodosForUser(userId, limit, offset)
}

func (s *TodoService) DeleteTodoByIDAndUserID(todoID string, userID uint) error {
	return s.TodoRepo.DeleteTodoByIDAndUserID(todoID, userID)
}

func (s *TodoService) UpdateTodoByIDAndUserID(todoID string, userID uint, updatedTodo *models.Todo) error {
	return s.TodoRepo.UpdateTodoByIDAndUserID(todoID, userID, updatedTodo)
}
