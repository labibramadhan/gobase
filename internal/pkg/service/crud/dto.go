package crud

import "errors"

// ErrNotFound is returned when an entity is not found in the database.
var ErrNotFound = errors.New("entity not found")

type PaginationResult struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
	TotalRows  int64 `json:"totalRows"`
	HasNext    bool  `json:"hasNext"`
}

// PageResult represents a paginated list of items.
type PageResult[T any] struct {
	Items      []T              `json:"items"`
	Pagination PaginationResult `json:"pagination"`
}

// FilterOperator defines the supported filter operators
type FilterOperator string

const (
	OperatorEqual              FilterOperator = "eq"
	OperatorNotEqual           FilterOperator = "neq"
	OperatorGreaterThan        FilterOperator = "gt"
	OperatorGreaterThanOrEqual FilterOperator = "gte"
	OperatorLessThan           FilterOperator = "lt"
	OperatorLessThanOrEqual    FilterOperator = "lte"
	OperatorLike               FilterOperator = "like"
	OperatorILike              FilterOperator = "ilike"
	OperatorIn                 FilterOperator = "in"
	OperatorNotIn              FilterOperator = "nin"
	OperatorIsNull             FilterOperator = "isnull"
	OperatorIsNotNull          FilterOperator = "isnotnull"
)

// LogicalOperator defines the logical operators for grouping filters
type LogicalOperator string

const (
	LogicalAnd LogicalOperator = "AND"
	LogicalOr  LogicalOperator = "OR"
)

// Pagination defines the pagination parameters
type Pagination struct {
	Page      int  `json:"page"`
	PageSize  int  `json:"pageSize"`
	WithCount bool `json:"withCount"`
}

// Sort defines the sorting parameters
type Sort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

// Filter defines a single filter criterion
type Filter struct {
	Field    string         `json:"field"`
	Operator FilterOperator `json:"operator"`
	Value    any            `json:"value"`
}

// FilterGroup defines a group of filters with a logical operator
type FilterGroup struct {
	Operator LogicalOperator `json:"operator"`
	Filters  []any           `json:"filters"` // Can be Filter or FilterGroup
}

// QueryOptions holds all the query parameters
type QueryOptions struct {
	Pagination *Pagination  `json:"pagination,omitempty"`
	Sorts      []Sort       `json:"sorts,omitempty"`
	Filters    *FilterGroup `json:"filters,omitempty"`
}

// NewQueryOptions creates a new QueryOptions with default pagination.
func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		Pagination: &Pagination{
			Page:     1,
			PageSize: 10,
		},
	}
}

// WithPagination sets the pagination for the query.
func (q *QueryOptions) WithPagination(page, pageSize int) *QueryOptions {
	q.Pagination = &Pagination{Page: page, PageSize: pageSize}
	return q
}

// WithSort adds a sort criterion to the query.
func (q *QueryOptions) WithSort(field, direction string) *QueryOptions {
	q.Sorts = append(q.Sorts, Sort{Field: field, Direction: direction})
	return q
}

// WithFilter sets the filter for the query.
func (q *QueryOptions) WithFilter(filter *FilterGroup) *QueryOptions {
	q.Filters = filter
	return q
}
