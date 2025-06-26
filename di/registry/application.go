package registry

import (
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
		cleanup, err := a.TransportREST.Run()
		<-a.ListenTerminateSignal(cleanup)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to run REST application")
		}
	}()
}

func (a *Application) RunTransportGraphQL() {
	// Run REST server in a goroutine to prevent blocking
	go func() {
		log.Info().Msg("starting GraphQL server")
		cleanup, err := a.TransportGraphQL.Run()
		<-a.ListenTerminateSignal(cleanup)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to run GraphQL application")
		}
	}()
}

func (a *Application) ListenTerminateSignal(cleanup CleanupFunc) chan os.Signal {
	if a.terminationSignal == nil {
		a.terminationSignal = make(chan os.Signal, 1)
		signal.Notify(a.terminationSignal, syscall.SIGINT, syscall.SIGTERM)

		// Handle termination in a goroutine
		go func() {
			<-a.terminationSignal
			log.Info().Msg("Received termination signal")

			cleanup()

			log.Info().Msg("Application shutting down gracefully")
			os.Exit(0)
		}()
	}
	return a.terminationSignal
}
