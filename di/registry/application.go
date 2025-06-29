package registry

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

type Application struct {
	TransportREST      IApplicationTransportREST
	TransportGraphQL   IApplicationTransportGraphQL
	TransportWatermill IApplicationTransportWatermill
	Initialize         InitializerFunc

	terminationSignal chan os.Signal
}

func NewApplication(
	transportREST IApplicationTransportREST,
	transportGraphQL IApplicationTransportGraphQL,
	transportWatermill IApplicationTransportWatermill,
	initializer InitializerFunc,
) *Application {
	return &Application{
		TransportREST:      transportREST,
		TransportGraphQL:   transportGraphQL,
		TransportWatermill: transportWatermill,
		Initialize:         initializer,
	}
}

func (a *Application) RunTransportREST(ctx context.Context) {
	// Run REST server in a goroutine to prevent blocking
	go func() {
		log.Info().Msg("starting REST server")
		err := a.TransportREST.Run(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to run REST application")
		}
	}()
}

func (a *Application) RunTransportGraphQL(ctx context.Context) {
	// Run REST server in a goroutine to prevent blocking
	go func() {
		log.Info().Msg("starting GraphQL server")
		err := a.TransportGraphQL.Run(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to run GraphQL application")
		}
	}()
}

func (a *Application) RunWatermill(ctx context.Context) {
	go func() {
		log.Info().Msg("starting Watermill")
		a.TransportWatermill.Run(ctx)
	}()
}

func (a *Application) Run(cleanup CleanupFunc) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a.Initialize()

	a.RunTransportREST(ctx)
	a.RunTransportGraphQL(ctx)
	a.RunWatermill(ctx)

	<-ctx.Done()

	defer func() {
		log.Info().Msg("Calling cleanup function...")
		cleanup()
	}()

	log.Info().Msg("Application shutting down gracefully...")
}
