package masterdataentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type RelProductVariantProductAttribute struct {
	bun.BaseModel `bun:"table:rel_product_variant_product_attribute"`

	Id          uuid.UUID `bun:"id,pk,type:uuid" validate:"uuid,required"`
	ProductId   uuid.UUID `bun:"product_id,type:uuid" validate:"uuid,required"`
	VariantId   uuid.UUID `bun:"product_variant_id,type:uuid" validate:"uuid,required"`
	AttributeId uuid.UUID `bun:"product_attribute_id,type:uuid" validate:"uuid,required"`
	Value       string    `validate:"required"`

	Attribute *ProductAttribute `bun:"rel:belongs-to,join:product_attribute_id=id"`
	Variant   *ProductVariant   `bun:"rel:belongs-to,join:product_variant_id=id"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete"`
}
