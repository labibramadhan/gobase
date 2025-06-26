package masterdataentity

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ bun.BeforeAppendModelHook = (*Product)(nil)
var _ bun.BeforeUpdateHook = (*Product)(nil)
var _ bun.BeforeDeleteHook = (*Product)(nil)

type Product struct {
	bun.BaseModel `bun:"table:product"`

	Id          uuid.UUID `bun:"id,pk,type:uuid" validate:"uuid,required"`
	Name        string    `validate:"required"`
	Description string    `validate:"required"`

	Variants []*ProductVariant `bun:"rel:has-many,join:id=product_id"`

	Version   int
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete"`

	IgnoreVersionOnUpdate bool `bun:"-" json:"-"`
}

func (u *Product) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		u.Version = 0
	case *bun.UpdateQuery:
		u.Version++
	}
	return nil
}

func (u *Product) BeforeUpdate(ctx context.Context, query *bun.UpdateQuery) error {
	if !u.IgnoreVersionOnUpdate {
		query.Where("version = ?", u.Version)
	}
	u.Version++
	return nil
}

func (u *Product) BeforeDelete(ctx context.Context, query *bun.DeleteQuery) error {
	if !u.IgnoreVersionOnUpdate {
		query.Where("version = ?", u.Version)
	}
	u.Version++
	return nil
}
