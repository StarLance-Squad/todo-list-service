package mdw

import (
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PaginationInfo struct {
	Page   int
	Limit  int
	Offset int
}

func getPaginationLimit() int {
	limitStr := os.Getenv("PAGINATION_LIMIT")
	if limitStr == "" {
		return 10
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 10
	}

	return limit
}

func PaginationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page == 0 {
			page = 1
		}
		limit := getPaginationLimit()

		pagination := &PaginationInfo{
			Page:   page,
			Limit:  limit,
			Offset: (page - 1) * limit,
		}

		c.Set("pagination", pagination)
		return next(c)
	}
}
