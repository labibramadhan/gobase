package productresolver

import (
	"context"

	productdto "gobase/internal/domain/product/dto"
)

func (r *ResolverModule) Create(ctx context.Context, input productdto.CreateProductInput) (*productdto.Product, error) {
	return r.productUseCase.Create(ctx, input)
}

func (r *ResolverModule) CreateAttribute(ctx context.Context, input productdto.CreateProductAttributeInput) (*productdto.ProductAttribute, error) {
	return r.productUseCase.CreateAttribute(ctx, input)
}
