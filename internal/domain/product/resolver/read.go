package productresolver

import (
	"context"

	"github.com/google/uuid"

	productdto "gobase/internal/domain/product/dto"
)

func (r *ResolverModule) FindById(ctx context.Context, id uuid.UUID) (*productdto.Product, error) {
	return r.productUseCase.FindById(ctx, id)
}

func (r *ResolverModule) FindAll(ctx context.Context, qop *productdto.ProductQop) (*productdto.ProductList, error) {
	return r.productUseCase.FindAll(ctx, qop)
}
