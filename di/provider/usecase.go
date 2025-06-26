package provider

import (
	"github.com/google/wire"

	productusecase "gobase/internal/domain/product/usecase"
)

var UseCaseSet = wire.NewSet(
	wire.Struct(new(productusecase.UseCaseOpts), "*"),
	productusecase.NewUseCase,
)
