package services

import (
	"golang.org/x/crypto/bcrypt"
	"todo-list-service/internal/db"
	"todo-list-service/internal/models"
)

type UserService struct {
	UserRepo db.UserRepository
}

func (s *UserService) CreateUser(user *models.User) (*models.User, error) {
	return s.UserRepo.CreateUser(user)
}

func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByUsername(username)
	if err != nil {
		// User not found or other error
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		// Password does not match
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateLastLogin(user *models.User) error {
	return s.UserRepo.UpdateLastLogin(user)
}

func (s *UserService) GetAllUsers(limit int, offset int) ([]models.User, error) {
	return s.UserRepo.GetAllUsers(limit, offset)
}

func (s *UserService) GetUsersCount() (int64, error) { return s.UserRepo.GetUsersCount() }
