package productrepository

import (
	"context"

	"github.com/uptrace/bun"

	masterdataentity "gobase/internal/db/masterdata/entity"
	"gobase/internal/pkg/service/buncrud"
)

// Repository defines the data access layer for the entire Product aggregate.
type Repository interface {
	// WithTx returns a new repository instance that uses the provided transaction.
	WithTx(ctx context.Context, tx bun.Tx) Repository

	// Product returns a repository scoped to the Product entity.
	Product() buncrud.BaseRepository[masterdataentity.Product]

	// Variants returns a repository scoped to the ProductVariant entity.
	Variant() buncrud.BaseRepository[masterdataentity.ProductVariant]

	// Attributes returns a repository scoped to the ProductAttribute entity.
	Attribute() buncrud.BaseRepository[masterdataentity.ProductAttribute]

	// VariantAttributeValue returns a repository for the variant-attribute relationship.
	VariantAttributeValue() buncrud.BaseRepository[masterdataentity.RelProductVariantProductAttribute]
}

// RepositoryModule is the implementation of the Repository interface.
type RepositoryModule struct {
	productsRepo        buncrud.BaseRepository[masterdataentity.Product]
	variantsRepo        buncrud.BaseRepository[masterdataentity.ProductVariant]
	attributesRepo      buncrud.BaseRepository[masterdataentity.ProductAttribute]
	attributeValuesRepo buncrud.BaseRepository[masterdataentity.RelProductVariantProductAttribute]
	db                  bun.IDB // Can be *bun.DB or *bun.Tx
}

type RepositoryOpts struct {
	Bun *bun.DB
}

// NewRepository creates a new repository for the Product aggregate.
func NewRepository(opts RepositoryOpts) Repository {
	return &RepositoryModule{
		productsRepo:        buncrud.NewBaseRepository[masterdataentity.Product](opts.Bun),
		variantsRepo:        buncrud.NewBaseRepository[masterdataentity.ProductVariant](opts.Bun),
		attributesRepo:      buncrud.NewBaseRepository[masterdataentity.ProductAttribute](opts.Bun),
		attributeValuesRepo: buncrud.NewBaseRepository[masterdataentity.RelProductVariantProductAttribute](opts.Bun),
		db:                  opts.Bun,
	}
}

// WithTx returns a new repository instance that uses the provided transaction.
func (r *RepositoryModule) WithTx(ctx context.Context, tx bun.Tx) Repository {
	return &RepositoryModule{
		productsRepo:        r.productsRepo.WithTx(ctx, tx),
		variantsRepo:        r.variantsRepo.WithTx(ctx, tx),
		attributesRepo:      r.attributesRepo.WithTx(ctx, tx),
		attributeValuesRepo: r.attributeValuesRepo.WithTx(ctx, tx),
		db:                  tx,
	}
}

func (r *RepositoryModule) Product() buncrud.BaseRepository[masterdataentity.Product] {
	return r.productsRepo
}

func (r *RepositoryModule) Variant() buncrud.BaseRepository[masterdataentity.ProductVariant] {
	return r.variantsRepo
}

func (r *RepositoryModule) Attribute() buncrud.BaseRepository[masterdataentity.ProductAttribute] {
	return r.attributesRepo
}

func (r *RepositoryModule) VariantAttributeValue() buncrud.BaseRepository[masterdataentity.RelProductVariantProductAttribute] {
	return r.attributeValuesRepo
}
