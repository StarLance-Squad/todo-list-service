package routes

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
	"todo-list-service/internal/auth"
	"todo-list-service/internal/db"
	"todo-list-service/internal/models"
)

func CreateTodo(c echo.Context) error {
	var todo models.Todo
	if err := c.Bind(&todo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := c.Get("user")
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token")
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unexpected token type")
	}

	// Inside CreateTodo function
	claims, ok := token.Claims.(*auth.JwtCustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "unexpected claims type")
	}

	todo.UserID = claims.UserID // Make sure this UserID is an integer

	// Save the new Todo
	result := db.DB.Create(&todo)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}

	return c.JSON(http.StatusCreated, todo)
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

func GetAllTodosByUserID(c echo.Context) error {
	// Get user from JWT token
	user := c.Get("user")
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token")
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unexpected token type")
	}

	claims, ok := token.Claims.(*auth.JwtCustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unexpected claims type")
	}

	// Print the JWT token data
	c.Logger().Print("JWT token data: ", claims)
	c.Logger().Print("User ID from JWT token: ", claims.UserID)

	userID := claims.UserID // Get user ID from JWT token

	// Get pagination parameters from query, with defaults
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 20
	}

	offset := (page - 1) * limit

	//// Query the database for todos by the user ID with pagination
	//var todos []models.Todo
	//result := db.DB.Preload("User").Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&todos)
	//if result.Error != nil {
	//	return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	//}

	// Inside GetAllTodosByUserID function
	var todos []TodoResponse
	result := db.DB.Model(&models.Todo{}).Select("id, title, description, completed, user_id, created_at, updated_at").Where("user_id = ?", userID).Limit(limit).Offset(offset).Scan(&todos)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}

	// Get the total count of todos for the user
	var count int64
	db.DB.Model(&models.Todo{}).Where("user_id = ?", userID).Count(&count)

	// Calculate the total number of pages
	totalPages := (count + int64(limit) - 1) / int64(limit)

	// Generate the next link if there are more pages
	var next string
	if page < int(totalPages) {
		next = fmt.Sprintf("/todos/get-all?page=%d&limit=%d", page+1, limit)
	}

	// Return the todos, count, and next link in the response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"todos": todos,
		"count": count,
		"next":  next,
	})
}
