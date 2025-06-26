package provider

import (
	"github.com/google/wire"

	"gobase/di/registry"
	productloader "gobase/internal/domain/product/dataloader"
)

var DataloaderGraphQLSet = wire.NewSet(
	productloader.NewDataloader,

	wire.Struct(new(registry.GraphQLDataloader), "*"),
)
