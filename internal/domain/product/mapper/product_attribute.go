package productmapper

import (
	masterdataentity "gobase/internal/db/masterdata/entity"
	productdto "gobase/internal/domain/product/dto"
)

func ProductAttributeValueEntityToDTO(productAttributeValueEntity *masterdataentity.RelProductVariantProductAttribute) *productdto.ProductAttributeValue {
	productAttributeValue := &productdto.ProductAttributeValue{
		ID:          productAttributeValueEntity.Id,
		Value:       productAttributeValueEntity.Value,
		AttributeId: productAttributeValueEntity.AttributeId,
	}

	if productAttributeValueEntity.Attribute != nil {
		productAttributeValue.Attribute = ProductAttributeEntityToDTO(productAttributeValueEntity.Attribute)
	}

	return productAttributeValue
}

func ProductAttributeEntityToDTO(productAttributeEntity *masterdataentity.ProductAttribute) *productdto.ProductAttribute {
	return &productdto.ProductAttribute{
		ID:   productAttributeEntity.Id,
		Name: productAttributeEntity.Name,
	}
}
