package masterdataentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProductAttribute struct {
	bun.BaseModel `bun:"table:product_attribute"`

	Id   uuid.UUID `bun:"id,pk,type:uuid" validate:"uuid,required"`
	Name string    `validate:"required"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete"`
}
