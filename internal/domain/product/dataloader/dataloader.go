package productloader

import (
	"github.com/google/uuid"
	dataloader "github.com/graph-gophers/dataloader/v7"

	masterdataentity "gobase/internal/db/masterdata/entity"
	productdto "gobase/internal/domain/product/dto"
	productmapper "gobase/internal/domain/product/mapper"
	productrepo "gobase/internal/domain/product/repository"
	"gobase/internal/pkg/service/gqldataloader"
)

// Dataloader holds all the dataloaders for the product domain.
type Dataloader struct {
	Variant        *dataloader.Loader[uuid.UUID, []*productdto.ProductVariant]
	AttributeValue *dataloader.Loader[uuid.UUID, []*productdto.ProductAttributeValue]
	Attribute      *dataloader.Loader[uuid.UUID, []*productdto.ProductAttribute]
}

// NewDataloader creates a new set of dataloaders for the product domain.
func NewDataloader(productRepo productrepo.Repository) *Dataloader {
	return &Dataloader{
		Variant:        dataloader.NewBatchedLoader(newVariantBatchFn(productRepo)),
		AttributeValue: dataloader.NewBatchedLoader(newAttributeValueBatchFn(productRepo)),
		Attribute:      dataloader.NewBatchedLoader(newAttributeBatchFn(productRepo)),
	}
}

// newVariantBatchFn creates a batch function for loading product variants using the generic batch function.
func newVariantBatchFn(repo productrepo.Repository) dataloader.BatchFunc[uuid.UUID, []*productdto.ProductVariant] {
	return gqldataloader.NewGenericBatchFn(
		repo.Variant(),
		[]string{"product_id"},
		func(item *masterdataentity.ProductVariant) uuid.UUID {
			return item.ProductId
		},
		nil,
		productmapper.ProductVariantEntityToDTO,
	)
}

// newAttributeValueBatchFn creates a batch function for loading product attribute values using the generic batch function.
func newAttributeValueBatchFn(repo productrepo.Repository) dataloader.BatchFunc[uuid.UUID, []*productdto.ProductAttributeValue] {
	return gqldataloader.NewGenericBatchFn(
		repo.VariantAttributeValue(),
		[]string{"product_variant_id"},
		func(item *masterdataentity.RelProductVariantProductAttribute) uuid.UUID {
			return item.VariantId
		},
		nil,
		productmapper.ProductAttributeValueEntityToDTO,
	)
}

// newAttributeBatchFn creates a batch function for loading product attribute values using the generic batch function.
func newAttributeBatchFn(repo productrepo.Repository) dataloader.BatchFunc[uuid.UUID, []*productdto.ProductAttribute] {
	return gqldataloader.NewGenericBatchFn(
		repo.Attribute(),
		[]string{"id"},
		func(item *masterdataentity.ProductAttribute) uuid.UUID {
			return item.Id
		},
		nil,
		productmapper.ProductAttributeEntityToDTO,
	)
}
