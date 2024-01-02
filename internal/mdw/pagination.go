package mdw

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

type PaginationInfo struct {
	Page   int
	Limit  int
	Offset int
}

func PaginationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page == 0 {
			page = 1
		}
		limit := 10 // todo - make this configurable

		pagination := &PaginationInfo{
			Page:   page,
			Limit:  limit,
			Offset: (page - 1) * limit,
		}

		c.Set("pagination", pagination)
		return next(c)
	}
}
