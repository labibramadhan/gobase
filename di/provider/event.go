package provider

import (
	"github.com/google/wire"

	producteventpublisher "gobase/internal/domain/product/event/publisher"
	producteventsubscriber "gobase/internal/domain/product/event/subscriber"
)

var EventSet = wire.NewSet(
	wire.Struct(new(producteventsubscriber.EventOpts), "*"),
	producteventsubscriber.NewEvent,
	wire.Struct(new(producteventpublisher.EventOpts), "*"),
	producteventpublisher.NewEvent,
)
