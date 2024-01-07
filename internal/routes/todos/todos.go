package todos

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"todo-list-service/internal/db"
	"todo-list-service/internal/models"
	"todo-list-service/internal/services"
)

func CreateTodo(c echo.Context, todoService *services.TodoService) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := claims["userId"].(float64)

	var todo models.Todo
	if err := c.Bind(&todo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	todo.UserID = userId

	// Save the new TodoItem
	newTodoItem, resultError := todoService.CreateTodo(&todo)
	if resultError != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, resultError.Error())
	}

	// Convert models.Todo to TodoResponse
	todoResponse := db.TodoResponse{
		ID:          newTodoItem.ID,
		Title:       newTodoItem.Title,
		Description: newTodoItem.Description,
		Completed:   newTodoItem.Completed,
		UserID:      newTodoItem.UserID,
	}

	return c.JSON(http.StatusCreated, todoResponse)
}

func GetAllTodosByUserID(c echo.Context, todoService *services.TodoService) error {

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := claims["userId"].(float64)

	// Print the JWT token data
	c.Logger().Print("JWT token data: ", claims)
	c.Logger().Print("userId: ", userId)

	// Get pagination parameters from query, with defaults
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page == 0 {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit

	//var todos []TodoResponse
	todos, resultError := todoService.GetAllTodosForUser(userId, limit, offset)
	if resultError != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, resultError.Error())
	}

	// Get the total count of todos for the user
	count, err := todoService.GetTodosCount(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	c.Logger().Print("count: ", count)

	// Return an empty list if there are no todos. Means no todos for the user
	if count == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"todos": []models.Todo{},
			"count": count,
			"next":  nil,
			"prev":  nil,
		})
	}

	// Calculate the total number of pages
	totalPages := (count + int64(limit) - 1) / int64(limit)
	c.Logger().Print("totalPages: ", totalPages)

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

	var getNext = func() interface{} {
		if next != "" {
			return next
		} else {
			return nil
		}
	}

	var getPrev = func() interface{} {
		if prev != "" {
			return prev
		} else {
			return nil
		}
	}

	// Return the todos, count, and next link in the response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"todos": todos,
		"count": count,
		"next":  getNext(),
		"prev":  getPrev(),
	})
}

func DeleteTodoByIDAndUserID(c echo.Context, todoService *services.TodoService) error {
	// Extracts the todoID and userID from the request parameters
	todoID := c.Param("todoID")
	userIDStr := c.Param("userID")

	// Parses the userID string to an unsigned integer
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid userID"})
	}

	// Calls the TodoService to delete a specific todo belonging to the user
	err = todoService.DeleteTodoByIDAndUserID(todoID, uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Returns a success message upon successful deletion
	return c.JSON(http.StatusOK, map[string]string{"message": "Todo deleted successfully"})
}

func UpdateTodoByIDAndUserID(c echo.Context, todoService *services.TodoService) error {
	// Extracts the todoID and userID from the request parameters
	todoID := c.Param("todoID")
	userIDStr := c.Param("userID")

	// Parses the userID string to an unsigned integer
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid userID"})
	}

	// Binds the request body to the updatedTodo model
	var updatedTodo models.Todo
	if err := c.Bind(&updatedTodo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Calls the TodoService to update a specific todo belonging to the user
	err = todoService.UpdateTodoByIDAndUserID(todoID, uint(userID), &updatedTodo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Returns a success message upon successful update
	return c.JSON(http.StatusOK, map[string]string{"message": "Todo updated successfully"})
}
