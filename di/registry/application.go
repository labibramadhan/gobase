package registry

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

type Application struct {
	TransportREST    IApplicationTransportREST
	TransportGraphQL IApplicationTransportGraphQL
	Initialize       InitializerFunc

	terminationSignal chan os.Signal
}

func NewApplication(
	transportREST IApplicationTransportREST,
	transportGraphQL IApplicationTransportGraphQL,
	initializer InitializerFunc,
) *Application {
	return &Application{
		TransportREST:    transportREST,
		TransportGraphQL: transportGraphQL,
		Initialize:       initializer,
	}
}

func (a *Application) RunTransportREST() {
	// Run REST server in a goroutine to prevent blocking
	go func() {
		log.Info().Msg("starting REST server")
		err := a.TransportREST.Run()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to run REST application")
		}
	}()
}

func (a *Application) RunTransportGraphQL() {
	// Run REST server in a goroutine to prevent blocking
	go func() {
		log.Info().Msg("starting GraphQL server")
		err := a.TransportGraphQL.Run()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to run GraphQL application")
		}
	}()
}

func (a *Application) Run(cleanup CleanupFunc) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a.Initialize()

	a.RunTransportREST()
	a.RunTransportGraphQL()

	<-ctx.Done()

	defer func() {
		log.Info().Msg("Calling cleanup function...")
		cleanup()
	}()

	log.Info().Msg("Application shutting down gracefully...")
}
