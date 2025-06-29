package productusecase

import (
	"context"

	"github.com/samber/lo"
	"github.com/uptrace/bun"

	masterdataentity "gobase/internal/db/masterdata/entity"
	productdto "gobase/internal/domain/product/dto"
	productmapper "gobase/internal/domain/product/mapper"
	"gobase/internal/pkg/service/otelsvc"
)

func (m *UseCaseModule) Create(ctx context.Context, productInput productdto.CreateProductInput) (*productdto.Product, error) {
	ctx, span := otelsvc.StartSpan(ctx, "ProductUseCase/Create")
	defer span.End()

	var err error

	err = m.sp.TransformAndValidateByTag(ctx, &productInput)
	if err != nil {
		return nil, err
	}

	productEntity := productInput.ToEntity(true)

	var createdProduct *masterdataentity.Product

	// The transaction will handle the creation of the product and all its related entities.
	err = m.bun.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		txRepo := m.repository.WithTx(ctx, tx)

		createdProduct, err = txRepo.Product().Create(ctx, productEntity)
		if err != nil {
			return err
		}
		err = m.productEventPublisher.PublishProductCreated(ctx, tx, createdProduct)
		if err != nil {
			return err
		}

		if len(productEntity.Variants) == 0 {
			return nil
		}

		_, err = txRepo.Variant().CreateBulk(ctx, productEntity.Variants)
		if err != nil {
			return err
		}

		attributeValues := lo.FlatMap(productEntity.Variants, func(variant *masterdataentity.ProductVariant, _ int) []*masterdataentity.RelProductVariantProductAttribute {
			return variant.Attributes
		})

		// Only create attribute values if there are any
		if len(attributeValues) > 0 {
			_, err = txRepo.VariantAttributeValue().CreateBulk(ctx, attributeValues)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return productmapper.ProductEntityToDTO(createdProduct), nil
}

func (m *UseCaseModule) CreateAttribute(ctx context.Context, input productdto.CreateProductAttributeInput) (*productdto.ProductAttribute, error) {
	ctx, span := otelsvc.StartSpan(ctx, "ProductUseCase/CreateAttribute")
	defer span.End()

	var err error

	err = m.sp.TransformAndValidateByTag(ctx, &input)
	if err != nil {
		return nil, err
	}

	attributeEntity := input.ToEntity(true)

	createdAttribute, err := m.repository.Attribute().Create(ctx, attributeEntity)
	if err != nil {
		return nil, err
	}

	return productmapper.ProductAttributeEntityToDTO(createdAttribute), nil
}
