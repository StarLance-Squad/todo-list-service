package auth

import (
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
	"time"
	"todo-list-service/internal/models"
	"todo-list-service/internal/services"
)

type (
	UserLogin struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	RegisterUser struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

// LoginHandler Login - Generate the token
func LoginHandler(c echo.Context, userService *services.UserService) (err error) {
	u := new(UserLogin)
	if err = c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(u); err != nil {
		return err
	}

	// Use userService to validate the user and get user details
	user, err := userService.AuthenticateUser(u.Username, u.Password)
	if err != nil {
		// Return a custom error response
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// Initialize a new instance of `Claims`
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.ID
	claims["username"] = user.Username
	claims["admin"] = false
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time is 24 hours

	// Signing token
	jwtSecret := os.Getenv("JWT_SECRET")
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return err
	}

	// Optionally, update last login using userService
	if err := userService.UpdateLastLogin(user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error on updating last login date-time.")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func RegisterUserHandler(c echo.Context, userService *services.UserService) (err error) {
	u := new(RegisterUser)
	if err = c.Bind(u); err != nil {
		c.Logger().Errorf("Error binding to RegisterUser: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(u); err != nil {
		c.Logger().Errorf("Error validating RegisterUser: %v", err)
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := &models.User{
		Username: u.Username,
		Password: string(hashedPassword),
		Email:    u.Email,
	}

	newCreatedUser, resultErr := userService.CreateUser(newUser)
	if resultErr != nil {
		if resultErr != nil {
			if strings.Contains(resultErr.Error(), "duplicate key value violates unique constraint") {
				return c.JSON(http.StatusConflict, map[string]string{"message": "Username already exists"})
			}
			return echo.NewHTTPError(http.StatusInternalServerError, resultErr.Error())
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"userId":    newCreatedUser.ID,
		"username":  newCreatedUser.Username,
		"email":     newCreatedUser.Email,
		"createdAt": newCreatedUser.CreatedAt,
		"updatedAt": newCreatedUser.UpdatedAt,
	})
}

func WhoIamHandler(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := claims["userId"].(float64)
	username := claims["username"].(string)
	admin := claims["admin"].(bool)
	exp := claims["exp"].(float64) // JWT exp claim is typically a Unix timestamp in seconds, decoded as float64
	return c.JSON(http.StatusOK, map[string]interface{}{
		"userId":   userId,
		"username": username,
		"admin":    admin,
		"exp":      exp,
	})
}

// todo - build a refresh token endpoint if needed

// todo - build a reset password endpoint
