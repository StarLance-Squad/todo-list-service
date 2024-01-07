package auth

import (
	"github.com/labstack/echo/v4"
	"todo-list-service/internal/services"
)

func InitRoutes(e *echo.Echo, userService *services.UserService) {
	authentication := e.Group("/authentication")
	authentication.POST("/login", func(c echo.Context) error { return LoginHandler(c, userService) })
	authentication.POST("/register", func(c echo.Context) error { return RegisterUserHandler(c, userService) })
	authentication.GET("/whoiam", WhoIamHandler)
}
