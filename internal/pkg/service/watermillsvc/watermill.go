package watermillsvc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	wsql "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/uptrace/bun"

	"gobase/internal/pkg/service/otelsvc"
)

// Service defines the Watermill service interface.
type Service interface {
	Run(ctx context.Context) error
	Shutdown() error
	AddSubscription(topic string, handlerFunc message.NoPublishHandlerFunc)
	Publish(ctx context.Context, topic string, messages ...*message.Message) error
	WithTx(tx wsql.ContextExecutor) (message.Publisher, error)
}

// ServiceOpts holds the options for the Watermill service.
type ServiceOpts struct {
	DB              wsql.Beginner
	Subscriber      message.Subscriber
	Publisher       message.Publisher // This is the "real" publisher (e.g., Kafka, RabbitMQ)
	Logger          watermill.LoggerAdapter
	BunSchemaConfig BunPostgreSQLSchemaConfig
}

// ServiceModule is the implementation of the Watermill service.
type ServiceModule struct {
	opts              ServiceOpts
	router            *message.Router
	publisher         message.Publisher    // This is the outbox publisher (always SQL-based)
	externalPublisher message.Publisher    // This is the external publisher (GoChannel, RabbitMQ, etc.)
	subscriber        message.Subscriber   // This is the external subscriber
	outboxSubscriber  *wsql.Subscriber     // This is the SQL subscriber for the outbox
	schemaAdapter     *BunPostgreSQLSchema // Schema adapter for the outbox
}

// NewService creates a new Watermill service with a multi-table outbox.
func NewService(opts ServiceOpts) (Service, error) {
	if opts.Logger == nil {
		opts.Logger = watermill.NewStdLogger(false, false)
	}

	router, err := message.NewRouter(message.RouterConfig{}, opts.Logger)
	if err != nil {
		return nil, err
	}

	router.AddMiddleware(
		middleware.CorrelationID,
		middleware.Retry{
			MaxRetries:      5,
			InitialInterval: time.Millisecond * 200,
			Logger:          opts.Logger,
		}.Middleware,
		middleware.Recoverer,
	)

	svc := &ServiceModule{
		opts:              opts,
		router:            router,
		subscriber:        opts.Subscriber, // External subscriber (GoChannel, RabbitMQ, etc.)
		externalPublisher: opts.Publisher,  // External publisher (GoChannel, RabbitMQ, etc.)
	}

	if opts.DB == nil {
		return nil, errors.New("database connection is required for outbox")
	}

	// Always create the schema adapter for the outbox pattern
	schemaAdapter, err := NewBunPostgreSQLSchema(opts.BunSchemaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create bun schema adapter: %w", err)
	}
	svc.schemaAdapter = schemaAdapter

	// Always create the SQL publisher for the outbox pattern
	outboxPublisher, err := wsql.NewPublisher(
		opts.DB,
		wsql.PublisherConfig{
			SchemaAdapter:        schemaAdapter,
			AutoInitializeSchema: true,
		},
		opts.Logger,
	)
	if err != nil {
		return nil, err
	}

	// The SQL publisher will initialize tables when publishing messages
	// We'll log the topics that will be used
	for topic := range opts.BunSchemaConfig.TopicToTableMap {
		opts.Logger.Info("Configured outbox table for topic", watermill.LogFields{"topic": topic})
	}

	svc.publisher = outboxPublisher

	outboxSubscriber, err := wsql.NewSubscriber(opts.DB, wsql.SubscriberConfig{
		SchemaAdapter:    schemaAdapter,
		OffsetsAdapter:   wsql.DefaultPostgreSQLOffsetsAdapter{},
		InitializeSchema: true,
		PollInterval:     time.Millisecond * 100,
	}, opts.Logger)
	if err != nil {
		return nil, err
	}
	svc.outboxSubscriber = outboxSubscriber

	// Register a forwarding handler for each topic in the topic-to-table mapping
	for topic, tableName := range opts.BunSchemaConfig.TopicToTableMap {
		handlerName := fmt.Sprintf("forwarder_handler_%s", tableName)

		// Initialize the outbox subscriber for this topic
		if err := outboxSubscriber.SubscribeInitialize(topic); err != nil {
			return nil, fmt.Errorf("failed to initialize outbox schema for topic %s: %w", topic, err)
		}

		opts.Logger.Info("Adding handler", watermill.LogFields{
			"handler": handlerName,
			"topic":   topic,
			"table":   tableName,
		})

		router.AddNoPublisherHandler(
			handlerName,
			topic, // Subscribe to the actual topic, not a forwarder topic
			outboxSubscriber,
			svc.forwardingHandler,
		)
	}

	return svc, nil
}

// forwardingHandler is a custom handler that forwards messages from the outbox to the external publisher.
func (s *ServiceModule) forwardingHandler(msg *message.Message) error {
	destinationTopic, ok := msg.Metadata[DestinationTopicKey]
	if !ok {
		s.opts.Logger.Error(
			"Message is missing destination topic, acking to avoid poison pill",
			nil,
			watermill.LogFields{"msg_uuid": msg.UUID},
		)
		return nil // Ack the message to prevent it from blocking the queue
	}

	// The message from the outbox is complete, forward it to the external publisher.
	// The router's middleware will handle acking/nacking.
	return s.externalPublisher.Publish(destinationTopic, msg)
}

// Run starts the Watermill router.
func (s *ServiceModule) Run(ctx context.Context) error {
	return s.router.Run(ctx)
}

// Shutdown gracefully stops the Watermill router.
func (s *ServiceModule) Shutdown() error {
	return s.router.Close()
}

// AddSubscription adds a new subscription handler to the router.
func (s *ServiceModule) AddSubscription(topic string, handlerFunc message.NoPublishHandlerFunc) {
	// Initialize the outbox subscriber for this topic
	if err := s.outboxSubscriber.SubscribeInitialize(topic); err != nil {
		s.opts.Logger.Error("failed to initialize outbox schema for topic", err, watermill.LogFields{"topic": topic})
		panic(err)
	}

	decoratedHandlerFunc := func(msg *message.Message) error {
		if msg != nil {
			_, span := otelsvc.StartSpanWithAttributes(msg.Context(), "Watermillsvc/Subscribe", map[string]any{
				"topic": topic,
			})
			defer span.End()
		}
		return handlerFunc(msg)
	}

	// Add the handler for the external subscriber (GoChannel, RabbitMQ, etc.)
	s.router.AddNoPublisherHandler(
		"handler_"+topic,
		topic,
		s.subscriber,
		decoratedHandlerFunc,
	)
}

// Publish publishes messages to the appropriate outbox table based on the topic mapping.
// This always uses the SQL outbox publisher, not the external publisher directly.
func (s *ServiceModule) Publish(ctx context.Context, topic string, messages ...*message.Message) error {
	ctx, span := otelsvc.StartSpanWithAttributes(ctx, "Watermillsvc/Publish", map[string]any{
		"topic": topic,
	})
	defer span.End()

	// For each message, add the destination topic to the metadata
	for _, msg := range messages {
		if msg.Metadata == nil {
			msg.Metadata = make(message.Metadata)
		}
		msg.Metadata[DestinationTopicKey] = topic
	}

	// Publish to the outbox using the original topic
	// The schema adapter will map it to the correct table internally
	return s.publisher.Publish(topic, messages...)
}

// WithTx returns a transactional publisher that writes messages to the outbox table.
func (s *ServiceModule) WithTx(tx wsql.ContextExecutor) (message.Publisher, error) {
	var executor wsql.ContextExecutor = tx

	if bunTx, ok := tx.(*bun.Tx); ok {
		executor = bunTx.Tx
	} else if bunTx, ok := tx.(bun.Tx); ok {
		executor = bunTx.Tx
	}

	return wsql.NewPublisher(
		executor,
		wsql.PublisherConfig{
			SchemaAdapter: s.schemaAdapter,
		},
		s.opts.Logger,
	)
}
