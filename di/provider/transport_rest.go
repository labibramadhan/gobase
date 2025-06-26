package provider

import (
	"github.com/google/wire"

	"gobase/di/registry"
	transportrest "gobase/transport/rest"
)

var TransportRESTDependencySet = wire.NewSet(
	wire.Struct(new(registry.RESTRouter), "*"),
)

var TransportRESTSet = wire.NewSet(
	TransportRESTDependencySet,
	transportrest.NewTransport,
)
