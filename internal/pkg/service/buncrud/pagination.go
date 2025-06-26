package buncrud

import (
	"github.com/uptrace/bun"

	"gobase/internal/pkg/service/crud"
)

// ApplyPagination applies pagination to the query
func ApplyPagination(query *bun.SelectQuery, pagination *crud.Pagination) *bun.SelectQuery {
	if pagination != nil && pagination.Page > 0 && pagination.PageSize > 0 {
		query.Limit(pagination.PageSize).Offset((pagination.Page - 1) * pagination.PageSize)
	}
	return query
}
