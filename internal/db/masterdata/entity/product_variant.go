package masterdataentity

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ bun.BeforeAppendModelHook = (*ProductVariant)(nil)
var _ bun.BeforeUpdateHook = (*ProductVariant)(nil)
var _ bun.BeforeDeleteHook = (*ProductVariant)(nil)

type ProductVariant struct {
	bun.BaseModel `bun:"table:product_variant"`

	Id              uuid.UUID `bun:"id,pk,type:uuid" validate:"uuid,required"`
	ProductId       uuid.UUID `bun:"product_id,type:uuid" validate:"uuid,required"`
	Name            string    `validate:"required"`
	SKU             string    `validate:"required"`
	Price           float64   `validate:"required"`
	DiscountedPrice float64   `validate:"required"`

	Product    *Product                             `bun:"rel:belongs-to,join:product_id=id"`
	Attributes []*RelProductVariantProductAttribute `bun:"rel:has-many,join:id=product_variant_id"`

	Version   int
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete"`

	IgnoreVersionOnUpdate bool `bun:"-" json:"-"`
}

func (u *ProductVariant) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		u.Version = 0
	case *bun.UpdateQuery:
		u.Version++
	}
	return nil
}

func (u *ProductVariant) BeforeUpdate(ctx context.Context, query *bun.UpdateQuery) error {
	if !u.IgnoreVersionOnUpdate {
		query.Where("version = ?", u.Version)
	}
	u.Version++
	return nil
}

func (u *ProductVariant) BeforeDelete(ctx context.Context, query *bun.DeleteQuery) error {
	if !u.IgnoreVersionOnUpdate {
		query.Where("version = ?", u.Version)
	}
	u.Version++
	return nil
}
