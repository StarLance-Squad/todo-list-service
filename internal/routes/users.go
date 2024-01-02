package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
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
	// Get pagination parameters from query, with defaults
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page == 0 {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit

	// Get users slice
	users, err := userService.GetAllUsers(limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Get total count of users
	count, err := userService.GetUsersCount()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Calculate the total number of pages
	totalPages := (count + int64(limit) - 1) / int64(limit)

	// Check if requested page is beyond the total pages
	if int64(page) > totalPages {
		return echo.NewHTTPError(http.StatusNotFound, "Page not found")
	}

	// Generate the next link if there are more pages
	var next string
	if page < int(totalPages) {
		next = fmt.Sprintf("/users?page=%d", page+1)
	}

	// Generate the previous link if it's not the first page
	var prev string
	if page > 1 {
		prev = fmt.Sprintf("/users?page=%d", page-1)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": users,
		"count": count,
		"next":  next,
		"prev":  prev,
	})
}
