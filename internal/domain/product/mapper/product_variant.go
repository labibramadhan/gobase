package productmapper

import (
	"github.com/samber/lo"

	masterdataentity "gobase/internal/db/masterdata/entity"
	productdto "gobase/internal/domain/product/dto"
)

func ProductVariantEntityToDTO(productVariantEntity *masterdataentity.ProductVariant) *productdto.ProductVariant {
	productVariant := &productdto.ProductVariant{
		ID:              productVariantEntity.Id,
		ProductID:       productVariantEntity.ProductId,
		Sku:             productVariantEntity.SKU,
		Price:           productVariantEntity.Price,
		DiscountedPrice: productVariantEntity.DiscountedPrice,
		CreatedAt:       productVariantEntity.CreatedAt,
		UpdatedAt:       productVariantEntity.UpdatedAt,
	}

	if len(productVariantEntity.Attributes) > 0 {
		productVariant.Attributes = lo.Map(productVariantEntity.Attributes, func(productAttributeValueEntity *masterdataentity.RelProductVariantProductAttribute, _ int) *productdto.ProductAttributeValue {
			return ProductAttributeValueEntityToDTO(productAttributeValueEntity)
		})
	}

	return productVariant
}
