package routes

import (
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, jwtMiddleware echo.MiddlewareFunc) {
	// First touch page
	e.GET("/", Home)

	// Auth
	e.POST("/login", Login)

	// User routes
	e.POST("/users", CreateUser)
	e.GET("/users", GetUsers)

	// Todos routes
	e.POST("/todos", CreateTodo, jwtMiddleware)
	e.POST("/todos/get-all", GetAllTodosByUserID, jwtMiddleware)
}
