package routes

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"todo-list-service/internal/auth"
	"todo-list-service/internal/db"
	"todo-list-service/internal/models"
)

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Assuming you have a function to validate the user and get user details
	user, err := validateUser(username, password)
	if err != nil {
		// Handle invalid user credentials
		return echo.ErrUnauthorized
	}

	// Generate JWT token
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := auth.GenerateToken(user.Username, user.ID, jwtSecret) // user.ID should be an integer
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating token")
	}

	// Update last login
	if err := user.UpdateLastLogin(db.DB); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error on updating last login date-time.")
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func validateUser(username, password string) (models.User, error) {
	// Fetch user from the database
	user, err := db.GetUserByUsername(username)
	if err != nil {
		// User not found or other error
		return models.User{}, err
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Password does not match
		return models.User{}, err
	}

	// User is valid
	return user, nil
}
