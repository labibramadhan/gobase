//go:build wireinject
// +build wireinject

package container

import (
	"github.com/google/wire"

	"gobase/di/provider"
	"gobase/di/registry"
)

func InitializeApplication(appContext registry.ApplicationContext) (*registry.Application, registry.CleanupFunc, error) {
	wire.Build(
		wire.FieldsOf(new(registry.ApplicationContext), "AgentListen", "ReadConfig"),
		provider.ProvideConfig,
		provider.Initializer,
		provider.InfrastructureSet,
		provider.TransportSet,
		provider.RepositorySet,
		provider.UseCaseSet,
		provider.EventSet,
		provider.DataloaderGraphQLSet,
		provider.MiddlewareGraphQLSet,
		provider.ServiceSet,
		provider.WatermillSet,
		registry.NewApplication,
	)
	return nil, nil, nil
}
