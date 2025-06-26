package provider

import (
	"github.com/google/wire"

	productrepository "gobase/internal/domain/product/repository"
)

var RepositorySet = wire.NewSet(
	wire.Struct(new(productrepository.RepositoryOpts), "*"),
	productrepository.NewRepository,
)
