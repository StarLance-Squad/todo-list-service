package users

import (
	"github.com/labstack/echo/v4"
	"todo-list-service/internal/services"
)

func InitRoutes(e *echo.Echo, userService *services.UserService) {
	users := e.Group("/users")
	users.GET("", func(c echo.Context) error { return GetUsers(c, userService) })
}
