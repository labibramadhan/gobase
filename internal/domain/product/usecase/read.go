package productusecase

import (
	"context"

	"github.com/google/uuid"

	productdto "gobase/internal/domain/product/dto"
	productmapper "gobase/internal/domain/product/mapper"
	"gobase/internal/pkg/service/crud"
	"gobase/internal/pkg/service/otelsvc"
)

func (m *UseCaseModule) FindById(ctx context.Context, id uuid.UUID) (*productdto.Product, error) {
	ctx, span := otelsvc.StartSpan(ctx, "ProductUseCase/FindById")
	defer span.End()

	productEntity, err := m.repository.Product().FindByID(ctx, id.String())
	if err != nil {
		return nil, err
	}

	return productmapper.ProductEntityToDTO(productEntity), nil
}

func (m *UseCaseModule) FindAll(ctx context.Context, qop *productdto.ProductQop) (*crud.PageResult[*productdto.Product], error) {
	ctx, span := otelsvc.StartSpan(ctx, "ProductUseCase/FindAll")
	defer span.End()

	var options *crud.QueryOptions

	if qop != nil {
		qop.WithAllowedSorts([]string{"id", "name", "created_at", "updated_at"})
		options = qop.ToQueryOptions()
	}

	// Get product entities from repository
	entityResult, err := m.repository.Product().FindAll(ctx, options)
	if err != nil {
		return nil, err
	}

	// Convert product entities to DTOs
	productDTOs := make([]*productdto.Product, len(entityResult.Items))
	for i, p := range entityResult.Items {
		productDTOs[i] = productmapper.ProductEntityToDTO(&p)
	}

	// Create the final DTO page result
	return &crud.PageResult[*productdto.Product]{
		Items:      productDTOs,
		Pagination: entityResult.Pagination,
	}, nil
}
