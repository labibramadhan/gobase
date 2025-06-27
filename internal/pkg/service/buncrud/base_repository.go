package buncrud

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/uptrace/bun"

	"gobase/internal/pkg/service/crud"
	"gobase/internal/pkg/service/otelsvc"
)

// BaseRepository defines the common repository interface
type BaseRepository[T any] interface {
	// WithTx returns a new repository instance that uses the provided transaction.
	WithTx(ctx context.Context, tx bun.Tx) BaseRepository[T]

	FindAll(ctx context.Context, options *crud.QueryOptions) (*crud.PageResult[T], error)
	FindIn(ctx context.Context, column string, values []any, options *crud.QueryOptions) ([]*T, error)
	FindByID(ctx context.Context, id string) (*T, error)
	Create(ctx context.Context, entity *T) (*T, error)
	CreateBulk(ctx context.Context, entities []*T) ([]*T, error)
	Update(ctx context.Context, entity *T) (*T, error)
	Delete(ctx context.Context, id string) error
	HardDelete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
	QueryBuilder(ctx context.Context, options *crud.QueryOptions) *bun.SelectQuery
}

// BaseRepositoryImpl implements BaseRepository
type BaseRepositoryImpl[T any] struct {
	db bun.IDB
}

// NewBaseRepository creates a new BaseRepository
func NewBaseRepository[T any](db bun.IDB) BaseRepository[T] {
	return &BaseRepositoryImpl[T]{db: db}
}

// WithTx returns a new repository instance that uses the provided transaction.
func (r *BaseRepositoryImpl[T]) WithTx(ctx context.Context, tx bun.Tx) BaseRepository[T] {
	_, span := otelsvc.StartSpan(ctx, "buncrud.WithTx")
	defer span.End()

	return &BaseRepositoryImpl[T]{
		db: tx,
	}
}

// QueryBuilder creates a new query builder with applied options
func (r *BaseRepositoryImpl[T]) QueryBuilder(ctx context.Context, options *crud.QueryOptions) *bun.SelectQuery {
	_, span := otelsvc.StartSpan(ctx, "buncrud.QueryBuilder")
	defer span.End()

	var entity T
	query := r.db.NewSelect().Model(&entity)

	opts := crud.NewQueryOptions()
	if options != nil {
		opts = options
	}

	// Apply filters
	if opts.Filters != nil {
		ApplyFilters(query, opts.Filters)
	}

	// Apply pagination
	if opts.Pagination != nil && opts.Pagination.Page > 0 && opts.Pagination.PageSize > 0 {
		if opts.Pagination.PageSize < 1 {
			opts.Pagination.PageSize = 10
		}
		ApplyPagination(query, opts.Pagination)
	}

	// Apply sorting
	for _, s := range opts.Sorts {
		direction := strings.ToUpper(s.Direction)
		if direction != "ASC" && direction != "DESC" {
			direction = "ASC"
		}
		query.Order(fmt.Sprintf("%s %s", s.Field, direction))
	}

	return query
}

// FindAll finds all entities matching the given options, with pagination and without count.
func (r *BaseRepositoryImpl[T]) FindAll(ctx context.Context, options *crud.QueryOptions) (*crud.PageResult[T], error) {
	ctx, span := otelsvc.StartSpan(ctx, "buncrud.FindAll")
	defer span.End()

	var entities []T
	var count int
	var err error

	// Build the base query
	query := r.QueryBuilder(ctx, options)

	// Get the total count of items if requested
	if options != nil && options.Pagination != nil && options.Pagination.WithCount {
		count, err = query.Count(ctx)
		if err != nil {
			return nil, err
		}
	}

	if options == nil {
		options = &crud.QueryOptions{
			Pagination: &crud.Pagination{
				Page:     1,
				PageSize: 10,
			},
		}
	}

	// Execute the query
	if err := query.Scan(ctx, &entities); err != nil {
		return nil, err
	}

	// Create the page result
	pageResult := &crud.PageResult[T]{
		Items: entities,
	}

	if options != nil && options.Pagination != nil {
		pageResult.Pagination.Page = options.Pagination.Page
		pageResult.Pagination.PageSize = options.Pagination.PageSize

		if options.Pagination.WithCount {
			pageResult.Pagination.TotalRows = int64(count)
			if options.Pagination.PageSize > 0 {
				pageResult.Pagination.TotalPages = int(math.Ceil(float64(count) / float64(options.Pagination.PageSize)))
			}
			pageResult.Pagination.HasNext = pageResult.Pagination.Page*pageResult.Pagination.PageSize < int(pageResult.Pagination.TotalRows)
		}
	}

	return pageResult, nil
}

// FindIn finds multiple entities where the given column is in the given values.
func (r *BaseRepositoryImpl[T]) FindIn(ctx context.Context, column string, values []any, options *crud.QueryOptions) ([]*T, error) {
	ctx, span := otelsvc.StartSpan(ctx, "buncrud.FindIn")
	defer span.End()

	var entities []T

	if len(values) == 0 {
		return []*T{},
			nil
	}

	query := r.QueryBuilder(ctx, options).Where(fmt.Sprintf("%s IN (?)", column), bun.In(values))

	if err := query.Scan(ctx, &entities); err != nil {
		return nil, err
	}

	result := make([]*T, len(entities))
	for i := range entities {
		result[i] = &entities[i]
	}

	return result, nil
}

// FindByID finds an entity by ID with optional relations.
// It returns ErrNotFound if the entity is not found.
func (r *BaseRepositoryImpl[T]) FindByID(ctx context.Context, id string) (*T, error) {
	ctx, span := otelsvc.StartSpan(ctx, "buncrud.FindByID")
	defer span.End()

	var entity T

	query := r.QueryBuilder(ctx, nil).Where("id = ?", id)

	if err := query.Scan(ctx, &entity); err != nil {
		if err == sql.ErrNoRows {
			return nil, crud.ErrNotFound
		}
		return nil, err
	}

	return &entity, nil
}

// Create creates a new entity and returns it.
func (r *BaseRepositoryImpl[T]) Create(ctx context.Context, entity *T) (*T, error) {
	ctx, span := otelsvc.StartSpan(ctx, "buncrud.Create")
	defer span.End()

	_, err := r.db.NewInsert().Model(entity).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// CreateBulk creates multiple entities in a single query.
func (r *BaseRepositoryImpl[T]) CreateBulk(ctx context.Context, entities []*T) ([]*T, error) {
	ctx, span := otelsvc.StartSpan(ctx, "buncrud.CreateBulk")
	defer span.End()

	if len(entities) == 0 {
		return entities, nil
	}
	_, err := r.db.NewInsert().Model(&entities).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

// Update updates an existing entity and returns it.
// It returns ErrNotFound if the entity does not exist.
func (r *BaseRepositoryImpl[T]) Update(ctx context.Context, entity *T) (*T, error) {
	ctx, span := otelsvc.StartSpan(ctx, "buncrud.Update")
	defer span.End()

	res, err := r.db.NewUpdate().Model(entity).WherePK().Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, crud.ErrNotFound
	}

	return entity, nil
}

// Delete performs a soft delete on an entity.
// It returns ErrNotFound if the entity does not exist.
func (r *BaseRepositoryImpl[T]) Delete(ctx context.Context, id string) error {
	ctx, span := otelsvc.StartSpan(ctx, "buncrud.Delete")
	defer span.End()

	var entity T
	res, err := r.db.NewUpdate().Model(&entity).
		Set("deleted_at = NOW()").
		Where("id = ?", id).
		Where("deleted_at IS NULL"). // Ensure we only soft-delete once
		Exec(ctx)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return crud.ErrNotFound
	}

	return nil
}

// HardDelete deletes an entity by ID.
// It returns ErrNotFound if the entity does not exist.
func (r *BaseRepositoryImpl[T]) HardDelete(ctx context.Context, id string) error {
	ctx, span := otelsvc.StartSpan(ctx, "buncrud.HardDelete")
	defer span.End()

	var entity T
	res, err := r.db.NewDelete().Model(&entity).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return crud.ErrNotFound
	}

	return nil
}

// Exists checks if an entity with the given ID exists
func (r *BaseRepositoryImpl[T]) Exists(ctx context.Context, id string) (bool, error) {
	ctx, span := otelsvc.StartSpan(ctx, "buncrud.Exists")
	defer span.End()

	var entity T
	exists, err := r.db.NewSelect().Model(&entity).
		Where("id = ?", id).
		Exists(ctx)
	return exists, err
}
