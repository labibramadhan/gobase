package provider

import (
	"github.com/google/wire"

	"gobase/di/registry"
	productresolver "gobase/internal/domain/product/resolver"
	transportgraphql "gobase/transport/graphql"
)

var TransportGraphQLDependencySet = wire.NewSet(
	wire.Struct(new(productresolver.ResolverOptions), "*"),
	productresolver.NewResolver,

	wire.Struct(new(registry.GraphQLResolver), "*"),
)

var TransportGraphQLSet = wire.NewSet(
	TransportGraphQLDependencySet,
	transportgraphql.NewTransport,
)
