package productdto

import (
	"time"

	"github.com/google/uuid"

	masterdataentity "gobase/internal/db/masterdata/entity"
)

type CreateProductInput struct {
	Name        string                      `json:"name" validate:"min=3"`
	Description string                      `json:"description"`
	Variants    []CreateProductVariantInput `json:"variants"`
}

func (c *CreateProductInput) ToEntity(isNew bool) *masterdataentity.Product {
	var productId = uuid.New()

	var product = &masterdataentity.Product{
		Id:          productId,
		Name:        c.Name,
		Description: c.Description,
	}

	if isNew {
		product.CreatedAt = time.Now()
		product.Version = 1
	} else {
		product.UpdatedAt = time.Now()
	}

	var variants []*masterdataentity.ProductVariant
	for _, v := range c.Variants {
		variantEntity := v.ToEntity(isNew)
		variantEntity.ProductId = productId

		for _, attr := range v.Attributes {
			attrEntity := attr.ToEntity(isNew)
			attrEntity.ProductId = productId
			attrEntity.VariantId = variantEntity.Id
			variantEntity.Attributes = append(variantEntity.Attributes, attrEntity)
		}
		variants = append(variants, variantEntity)
	}
	product.Variants = variants

	return product
}

type CreateProductVariantInput struct {
	Sku             string                             `json:"sku"`
	Price           float64                            `json:"price"`
	DiscountedPrice float64                            `json:"discounted_price"`
	Attributes      []CreateProductAttributeValueInput `json:"attributes"`
}

func (c *CreateProductVariantInput) ToEntity(isNew bool) *masterdataentity.ProductVariant {
	productVariant := &masterdataentity.ProductVariant{
		SKU:             c.Sku,
		Price:           c.Price,
		DiscountedPrice: c.DiscountedPrice,
	}

	if isNew {
		productVariant.Id = uuid.New()
		productVariant.CreatedAt = time.Now()
		productVariant.Version = 1
	} else {
		productVariant.UpdatedAt = time.Now()
	}

	return productVariant
}

type CreateProductAttributeValueInput struct {
	ID    uuid.UUID `json:"id"`
	Value string    `json:"value"`
}

func (c *CreateProductAttributeValueInput) ToEntity(isNew bool) *masterdataentity.RelProductVariantProductAttribute {
	productAttributeValue := &masterdataentity.RelProductVariantProductAttribute{
		AttributeId: c.ID,
		Value:       c.Value,
	}

	if isNew {
		productAttributeValue.Id = uuid.New()
		productAttributeValue.CreatedAt = time.Now()
	} else {
		productAttributeValue.UpdatedAt = time.Now()
	}

	return productAttributeValue
}

type CreateProductAttributeInput struct {
	Name string `json:"name"`
}

func (c *CreateProductAttributeInput) ToEntity(isNew bool) *masterdataentity.ProductAttribute {
	productAttribute := &masterdataentity.ProductAttribute{
		Name: c.Name,
	}

	if isNew {
		productAttribute.Id = uuid.New()
		productAttribute.CreatedAt = time.Now()
	} else {
		productAttribute.UpdatedAt = time.Now()
	}

	return productAttribute
}
