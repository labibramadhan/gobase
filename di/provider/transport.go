package provider

import (
	"github.com/google/wire"
)

var TransportSet = wire.NewSet(
	TransportRESTSet,
	TransportGraphQLSet,
)
