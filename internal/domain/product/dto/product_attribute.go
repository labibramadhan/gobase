package productdto

import (
	"github.com/google/uuid"
)

type ProductAttribute struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ProductAttributeValue struct {
	ID          uuid.UUID         `json:"id"`
	Value       string            `json:"value"`
	AttributeId uuid.UUID         `json:"attribute_id"`
	Attribute   *ProductAttribute `json:"attribute"`
}
