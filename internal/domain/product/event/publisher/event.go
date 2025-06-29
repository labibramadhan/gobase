package producteventpublisher

import (
	"context"

	wsql "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"

	masterdataentity "gobase/internal/db/masterdata/entity"
	"gobase/internal/pkg/service/watermillsvc"
)

type Event interface {
	PublishProductCreated(ctx context.Context, tx wsql.ContextExecutor, product *masterdataentity.Product) error
}

type EventModule struct {
	watermillsvc watermillsvc.Service
}

type EventOpts struct {
	Watermillsvc watermillsvc.Service
}

func NewEvent(opts EventOpts) Event {
	event := &EventModule{
		watermillsvc: opts.Watermillsvc,
	}

	return event
}

func (m *EventModule) PublishProductCreated(ctx context.Context, tx wsql.ContextExecutor, product *masterdataentity.Product) error {
	msg, err := watermillsvc.BuildNewMessage(product)
	if err != nil {
		return err
	}
	publisher, err := m.watermillsvc.WithTx(tx)
	if err != nil {
		return err
	}
	return publisher.Publish("product.created", msg)
}
