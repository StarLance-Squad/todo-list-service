package routes

import (
	"github.com/labstack/echo/v4"
	"todo-list-service/internal/services"
)

func Init(e *echo.Echo, svc *services.Services, jwtMiddleware echo.MiddlewareFunc) {

	// First touch page
	e.GET("/", Home)

	// Auth routes
	e.POST("/login", func(c echo.Context) error {
		return Login(c, svc.UserService)
	})

	// ---------------- User routes ---------------- //
	e.POST("/users", func(c echo.Context) error {
		return CreateUserHandler(c, svc.UserService)
	})

	// todo - provide a user role Admin, and only Admins can access this route
	e.GET("/users", func(c echo.Context) error {
		return GetUsers(c, svc.UserService)
	}, jwtMiddleware)
	// ---------------- User routes ---------------- //

	// ---------------- Todos routes --------------- //
	e.POST("/todos", func(c echo.Context) error {
		return CreateTodo(c, svc.TodoService)
	}, jwtMiddleware)

	// Note! Only auth user can access this route and get his own todos
	e.GET("/todos", func(c echo.Context) error {
		return GetAllTodosByUserID(c, svc.TodoService)
	}, jwtMiddleware)
	// ---------------- Todos routes --------------- //
}
