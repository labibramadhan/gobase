package provider

import (
	"github.com/google/wire"

	transportwatermill "gobase/transport/watermill"
)

var TransportWatermillSet = wire.NewSet(
	wire.Struct(new(transportwatermill.TransportOpts), "*"),
	transportwatermill.NewTransport,
)
