package services

import (
	"golang.org/x/crypto/bcrypt"
	"todo-list-service/internal/db"
	middleware "todo-list-service/internal/mdw"
	"todo-list-service/internal/models"
	"todo-list-service/internal/utils"
)

type UserService struct {
	UserRepo db.UserRepository
}

func (s *UserService) CreateUser(user *models.User) (*models.User, error) {
	return s.UserRepo.CreateUser(user)
}

func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByUsername(username)

	// if any error occurs, we return nil and the error for Invalid credentials message
	if err != nil {
		// User not found or DB error
		return nil, err
	}

	// Check if the user object is nil
	if user == nil {
		return nil, err
	}

	// Correctly capture error from password comparison
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// Password does not match
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateLastLogin(user *models.User) error {
	return s.UserRepo.UpdateLastLogin(user)
}

func (s *UserService) GetAllUsersWithPagination(pagination *middleware.PaginationInfo, basePath string) (*PaginatedResponse[models.User], error) {
	totalCount, err := s.UserRepo.GetUsersCount()
	if err != nil {
		return nil, err
	}

	users, err := s.UserRepo.GetAllUsers(pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}

	next, prev := utils.GeneratePaginationLinks(pagination, totalCount, basePath)

	return &PaginatedResponse[models.User]{
		Data:       users,
		TotalCount: totalCount,
		Next:       next,
		Prev:       prev,
	}, nil
}

func (s *UserService) GetUsersCount() (int64, error) { return s.UserRepo.GetUsersCount() }
