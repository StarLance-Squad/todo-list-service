package users

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	middleware "todo-list-service/internal/mdw"
	"todo-list-service/internal/services"
)

func GetUsers(c echo.Context, userService *services.UserService) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	isAdmin := claims["admin"].(bool)
	if !isAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "You are not allowed to access this route")
	}

	pagination := c.Get("pagination").(*middleware.PaginationInfo)

	response, err := userService.GetAllUsersWithPagination(pagination, "/users")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}
