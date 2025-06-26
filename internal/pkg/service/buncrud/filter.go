package buncrud

import (
	"encoding/json"
	"fmt"

	"github.com/uptrace/bun"

	"gobase/internal/pkg/service/crud"
)

// ApplyFilters applies the filter group to the query
func ApplyFilters(query *bun.SelectQuery, filterGroup *crud.FilterGroup) {
	if filterGroup == nil {
		return
	}

	// Build the WHERE clause for the group
	query.WhereGroup(string(filterGroup.Operator), func(q *bun.SelectQuery) *bun.SelectQuery {
		for _, f := range filterGroup.Filters {
			switch v := f.(type) {
			case crud.Filter:
				applyFilter(q, v)
			case crud.FilterGroup:
				ApplyFilters(q, &v)
			default:
				// Handle potential marshaling from map[string]interface{}
				if marshaled, err := json.Marshal(f); err == nil {
					var concreteFilter crud.Filter
					if err := json.Unmarshal(marshaled, &concreteFilter); err == nil {
						applyFilter(q, concreteFilter)
						continue
					}

					var concreteGroup crud.FilterGroup
					if err := json.Unmarshal(marshaled, &concreteGroup); err == nil {
						ApplyFilters(q, &concreteGroup)
						continue
					}
				}
			}
		}
		return q
	})
}

// applyFilter applies a single filter to the query
func applyFilter(q *bun.SelectQuery, filter crud.Filter) {
	switch filter.Operator {
	case crud.OperatorEqual:
		q.Where("? = ?", bun.Ident(filter.Field), filter.Value)
	case crud.OperatorNotEqual:
		q.Where("? != ?", bun.Ident(filter.Field), filter.Value)
	case crud.OperatorGreaterThan:
		q.Where("? > ?", bun.Ident(filter.Field), filter.Value)
	case crud.OperatorGreaterThanOrEqual:
		q.Where("? >= ?", bun.Ident(filter.Field), filter.Value)
	case crud.OperatorLessThan:
		q.Where("? < ?", bun.Ident(filter.Field), filter.Value)
	case crud.OperatorLessThanOrEqual:
		q.Where("? <= ?", bun.Ident(filter.Field), filter.Value)
	case crud.OperatorLike:
		q.Where("? LIKE ?", bun.Ident(filter.Field), fmt.Sprintf("%%%v%%", filter.Value))
	case crud.OperatorILike:
		q.Where("? ILIKE ?", bun.Ident(filter.Field), fmt.Sprintf("%%%v%%", filter.Value))
	case crud.OperatorIn:
		if values, ok := filter.Value.([]interface{}); ok {
			q.Where("? IN (?)", bun.Ident(filter.Field), bun.In(values))
		}
	case crud.OperatorNotIn:
		if values, ok := filter.Value.([]interface{}); ok {
			q.Where("? NOT IN (?)", bun.Ident(filter.Field), bun.In(values))
		}
	case crud.OperatorIsNull:
		q.Where("? IS NULL", bun.Ident(filter.Field))
	case crud.OperatorIsNotNull:
		q.Where("? IS NOT NULL", bun.Ident(filter.Field))
	}
}
