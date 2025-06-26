package productusecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	productdto "gobase/internal/domain/product/dto"
	productrepository "gobase/internal/domain/product/repository"
	"gobase/internal/pkg/service/crud"
	structprocessor "gobase/internal/pkg/service/structprocessor"
)

type UseCase interface {
	Create(ctx context.Context, productInput productdto.CreateProductInput) (*productdto.Product, error)
	FindById(ctx context.Context, id uuid.UUID) (*productdto.Product, error)
	CreateAttribute(ctx context.Context, input productdto.CreateProductAttributeInput) (*productdto.ProductAttribute, error)
	FindAll(ctx context.Context, qop *productdto.ProductQop) (*crud.PageResult[*productdto.Product], error)
}

type UseCaseModule struct {
	bun        *bun.DB
	repository productrepository.Repository
	sp         structprocessor.StructProcessorService
}

type UseCaseOpts struct {
	Bun        *bun.DB
	Repository productrepository.Repository
	SP         structprocessor.StructProcessorService
}

func NewUseCase(opts UseCaseOpts) UseCase {
	return &UseCaseModule{
		bun:        opts.Bun,
		repository: opts.Repository,
		sp:         opts.SP,
	}
}
