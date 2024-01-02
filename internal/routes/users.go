package routes

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	middleware "todo-list-service/internal/mdw"
	"todo-list-service/internal/models"
	"todo-list-service/internal/services"
)

func CreateUserHandler(c echo.Context, userService *services.UserService) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Create the user and save to the database
	newCreatedUser, resultErr := userService.CreateUser(user)
	if resultErr != nil {
		if strings.Contains(resultErr.Error(), "duplicate key value violates unique constraint") {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Username already exists"})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, resultErr.Error())
	}

	return c.JSON(http.StatusCreated, newCreatedUser)
}

func GetUsers(c echo.Context, userService *services.UserService) error {
	pagination := c.Get("pagination").(*middleware.PaginationInfo)

	response, err := userService.GetAllUsersWithPagination(pagination, "/users")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}
