package productdto

import (
	"time"

	"k8s.io/utils/strings/slices"

	"gobase/internal/pkg/service/crud"
)

// ProductQopFilter defines the specific, allowed filters for products.
// Tags are used to map these fields to the underlying database query.
// This struct is based on the ProductFilterInput from the GraphQL schema.
type ProductQopFilter struct {
	Name         *string    `filter:"field:name;operator:like"`
	CreatedAt    *time.Time `filter:"field:created_at;operator:eq"`
	UpdatedAt    *time.Time `filter:"field:updated_at;operator:eq"`
	CreatedAtGte *time.Time `filter:"field:created_at;operator:gte"`
	CreatedAtLte *time.Time `filter:"field:created_at;operator:lte"`
}

// ProductQop (Query Options Provider) is an opinionated struct for product queries.
// It exposes specific filtering, sorting, and pagination options, preventing overly complex
// or performantly dangerous queries from the client.
type ProductQop struct {
	crud.QueryOptions
	Filters ProductQopFilter `json:"filters"`
}

// ToQueryOptions converts the opinionated ProductQop to the generic crud.QueryOptions
// that the repository layer can understand. It uses reflection to parse the `filter` tags.
func (q *ProductQop) ToQueryOptions() *crud.QueryOptions {
	qOpts := q.QueryOptions
	qOpts.Filters = crud.BuildFilter(q.Filters)
	return &qOpts
}

// WithAllowedSorts validates and sets the sorting options, ensuring only
// whitelisted fields can be used for sorting.
func (q *ProductQop) WithAllowedSorts(allowedSorts []string) *ProductQop {
	var validatedSorts []crud.Sort
	for _, s := range q.QueryOptions.Sorts {
		if slices.Contains(allowedSorts, s.Field) {
			validatedSorts = append(validatedSorts, s)
		}
	}
	q.QueryOptions.Sorts = validatedSorts
	return q
}
