package utils

import (
	"fmt"
	middleware "todo-list-service/internal/mdw"
)

func GeneratePaginationLinks(pagination *middleware.PaginationInfo, totalCount int64, basePath string) (string, string) {
	totalPages := (totalCount + int64(pagination.Limit) - 1) / int64(pagination.Limit)

	var next, prev string
	if pagination.Page < int(totalPages) {
		next = fmt.Sprintf("%s?page=%d", basePath, pagination.Page+1)
	}
	if pagination.Page > 1 {
		prev = fmt.Sprintf("%s?page=%d", basePath, pagination.Page-1)
	}

	return next, prev
}
