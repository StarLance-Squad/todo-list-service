package todos

import (
	"github.com/labstack/echo/v4"
	"todo-list-service/internal/services"
)

func InitRoutes(e *echo.Echo, todoService *services.TodoService) {
	todos := e.Group("/todos")
	todos.GET("", func(c echo.Context) error { return GetAllTodosByUserID(c, todoService) })
	todos.POST("", func(c echo.Context) error { return CreateTodo(c, todoService) })
	todos.PUT("", func(c echo.Context) error { return UpdateTodoByIDAndUserID(c, todoService) })
	todos.DELETE("", func(c echo.Context) error { return DeleteTodoByIDAndUserID(c, todoService) })
}
