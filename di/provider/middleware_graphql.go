package provider

import (
	"github.com/google/wire"

	middlewaregraphql "gobase/internal/pkg/middleware/graphql"
)

var MiddlewareGraphQLSet = wire.NewSet(
	middlewaregraphql.NewDataloader,
	middlewaregraphql.NewOtel,
)
