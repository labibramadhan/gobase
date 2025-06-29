package producteventsubscriber

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/zerolog/log"

	masterdataentity "gobase/internal/db/masterdata/entity"
	productusecase "gobase/internal/domain/product/usecase"
	"gobase/internal/pkg/service/watermillsvc"
)

type Event interface {
	Subscribe()
	HandleProductCreated(msg *message.Message) error
}

type EventModule struct {
	watermillsvc   watermillsvc.Service
	productUseCase productusecase.UseCase
}

type EventOpts struct {
	Watermillsvc   watermillsvc.Service
	ProductUseCase productusecase.UseCase
}

func NewEvent(opts EventOpts) Event {
	return &EventModule{
		watermillsvc:   opts.Watermillsvc,
		productUseCase: opts.ProductUseCase,
	}
}

func (m *EventModule) Subscribe() {
	m.watermillsvc.AddSubscription("product.created", m.HandleProductCreated)
}

func (m *EventModule) HandleProductCreated(msg *message.Message) error {
	entity := &masterdataentity.Product{}
	if err := json.Unmarshal(msg.Payload, entity); err != nil {
		return err
	}
	log.Info().Any("entity", entity).Msg("consuming message from topic product.created")
	return nil
}
