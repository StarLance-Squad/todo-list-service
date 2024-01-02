package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"todo-list-service/internal/auth"
	"todo-list-service/internal/services"
)

func Login(c echo.Context, userService *services.UserService) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Use userService to validate the user and get user details
	user, err := userService.AuthenticateUser(username, password)
	if err != nil {
		// Handle invalid user credentials
		return echo.ErrUnauthorized
	}

	// Generate JWT token
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := auth.GenerateToken(user.Username, user.ID, jwtSecret) // Ensure user.ID is the expected type for GenerateToken
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating token")
	}

	// Optionally, update last login using userService
	if err := userService.UpdateLastLogin(user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error on updating last login date-time.")
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
