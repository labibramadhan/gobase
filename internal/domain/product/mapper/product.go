package productmapper

import (
	"github.com/samber/lo"

	masterdataentity "gobase/internal/db/masterdata/entity"
	productdto "gobase/internal/domain/product/dto"
)

func ProductEntityToDTO(productEntity *masterdataentity.Product) *productdto.Product {
	product := &productdto.Product{
		ID:          productEntity.Id,
		Name:        productEntity.Name,
		Description: productEntity.Description,
		CreatedAt:   productEntity.CreatedAt,
		UpdatedAt:   productEntity.UpdatedAt,
	}

	if len(productEntity.Variants) > 0 {
		product.Variants = lo.Map(productEntity.Variants, func(productVariantEntity *masterdataentity.ProductVariant, _ int) *productdto.ProductVariant {
			return ProductVariantEntityToDTO(productVariantEntity)
		})
	}

	return product
}
