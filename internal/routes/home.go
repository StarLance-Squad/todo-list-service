package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Home(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "work"})
}
