package provider

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/google/wire"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"

	"gobase/config"
	"gobase/internal/pkg/service/watermillsvc"
)

// ProvideWatermillLogger creates a Watermill logger adapter for zerolog.
func ProvideWatermillLogger() watermill.LoggerAdapter {
	return watermillsvc.NewZerologLoggerAdapter(log.Logger)
}

// ProvideWatermillPublisher creates a GoChannel-based message publisher for external messaging.
func ProvideWatermillPublisher(logger watermill.LoggerAdapter) (message.Publisher, error) {
	// Using GoChannel as the external publisher
	return gochannel.NewGoChannel(
		gochannel.Config{
			BlockPublishUntilSubscriberAck: true,
		},
		logger,
	), nil
}

// ProvideWatermillSubscriber creates a GoChannel-based message subscriber for external messaging.
func ProvideWatermillSubscriber(publisher message.Publisher, logger watermill.LoggerAdapter) (message.Subscriber, error) {
	// For GoChannel, we use the same instance for both publishing and subscribing
	return publisher.(message.Subscriber), nil
}

// ProvideWatermillService creates a new Watermill service.
func ProvideWatermillService(cfg *config.MainConfig, logger watermill.LoggerAdapter, db *bun.DB, publisher message.Publisher, subscriber message.Subscriber) (watermillsvc.Service, error) {
	return watermillsvc.NewService(watermillsvc.ServiceOpts{
		DB:         db.DB,
		Subscriber: subscriber,
		Publisher:  publisher,
		Logger:     logger,
		BunSchemaConfig: watermillsvc.BunPostgreSQLSchemaConfig{
			TableNames:      cfg.Watermill.Outbox.TableNames,
			TopicToTableMap: cfg.Watermill.Outbox.TopicToTableMap,
		},
	})
}

// WatermillSet is the Wire provider set for the Watermill service.
var WatermillSet = wire.NewSet(
	ProvideWatermillLogger,
	ProvideWatermillPublisher,
	ProvideWatermillSubscriber,
	ProvideWatermillService,
)
