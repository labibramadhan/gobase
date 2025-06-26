package gqldataloader

import (
	"context"
	"strings"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"

	"gobase/internal/pkg/service/buncrud"
)

// NewGenericBatchFn creates a generic batch function for a dataloader.
// K is the key type (e.g., uuid.UUID).
// V is the entity type (e.g., masterdataentity.ProductVariant).
// T is the GraphQL DTO type.
func NewGenericBatchFn[K comparable, V any, T any](
	repo buncrud.BaseRepository[V],
	columns []string, // Support for composite keys
	getKey func(item *V) K,
	getKeyValues func(key K) []any, // New: Extracts values from composite key. Can be nil for single keys.
	transform func(item *V) *T, // Transformation function from V to T
) dataloader.BatchFunc[K, []*T] {
	return func(ctx context.Context, keys []K) []*dataloader.Result[[]*T] {
		results := make([]*dataloader.Result[[]*T], len(keys))

		// Initialize results with empty slices.
		for i := range keys {
			results[i] = &dataloader.Result[[]*T]{Data: make([]*T, 0)}
		}

		var anyKeys []any
		if len(columns) > 1 && getKeyValues != nil {
			// For composite keys, map each key struct to a slice of its values.
			anyKeys = lo.Map(keys, func(key K, _ int) any {
				return getKeyValues(key) // bun expects a slice of slices for composite IN
			})
		} else {
			// For single keys, map each key to its value directly.
			anyKeys = lo.Map(keys, func(key K, _ int) any {
				return key
			})
		}

		var items []*V // Changed back to slice of pointers
		var err error

		// Fetch all items for the given keys.
		if len(columns) == 1 {
			items, err = repo.FindIn(ctx, columns[0], anyKeys, nil)
		} else {
			items, err = repo.FindIn(ctx, "("+strings.Join(columns, ", ")+")", anyKeys, nil)
		}

		if err != nil {
			for i := range keys {
				results[i] = &dataloader.Result[[]*T]{Error: err}
			}
			return results
		}

		// Group items by key.
		grouped := lo.GroupBy(items, func(item *V) K { // Changed to pointer
			return getKey(item)
		})

		// Distribute and transform the grouped items.
		for i, key := range keys {
			if group, ok := grouped[key]; ok {
				results[i] = &dataloader.Result[[]*T]{
					Data: lo.Map(group, func(item *V, _ int) *T {
						return transform(item)
					}),
				}
			}
		}

		return results
	}
}
