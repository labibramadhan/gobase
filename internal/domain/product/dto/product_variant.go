package productdto

import (
	"time"

	"github.com/google/uuid"
)

type ProductVariant struct {
	ID              uuid.UUID                `json:"id"`
	ProductID       uuid.UUID                `json:"product_id"`
	Sku             string                   `json:"sku"`
	Price           float64                  `json:"price"`
	DiscountedPrice float64                  `json:"discounted_price"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	Attributes      []*ProductAttributeValue `json:"attributes"`
}
