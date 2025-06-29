package transportwatermill

import (
	"context"

	"gobase/di/registry"
	producteventsubscriber "gobase/internal/domain/product/event/subscriber"
	"gobase/internal/pkg/service/watermillsvc"
)

type TransportModule struct {
	watermillService       watermillsvc.Service
	productEventSubscriber producteventsubscriber.Event
}

type TransportOpts struct {
	WatermillService       watermillsvc.Service
	ProductEventSubscriber producteventsubscriber.Event
}

func NewTransport(opts TransportOpts) (registry.IApplicationTransportWatermill, registry.CleanupFunc) {
	transportModule := &TransportModule{
		watermillService:       opts.WatermillService,
		productEventSubscriber: opts.ProductEventSubscriber,
	}

	return transportModule, func() {}
}

func (m *TransportModule) Run(ctx context.Context) error {
	m.productEventSubscriber.Subscribe()

	// Use the passed context instead of creating a new one
	// This ensures proper context propagation and shutdown handling
	return m.watermillService.Run(ctx)
}
