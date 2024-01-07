package routes

import (
	"github.com/labstack/echo/v4"
	"todo-list-service/internal/routes/auth"
	"todo-list-service/internal/routes/todos"
	"todo-list-service/internal/routes/users"
	"todo-list-service/internal/services"
)

func Init(e *echo.Echo, svc *services.Services) {
	auth.InitRoutes(e, svc.UserService)
	users.InitRoutes(e, svc.UserService)
	todos.InitRoutes(e, svc.TodoService)
}
