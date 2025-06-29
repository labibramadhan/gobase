package registry

import (
	"context"

	"github.com/google/gops/agent"

	productloader "gobase/internal/domain/product/dataloader"
	productresolver "gobase/internal/domain/product/resolver"
)

type CleanupFunc = func()
type InitializerFunc = func()

type AgentListen func(opts agent.Options) error
type ReadConfig func(cfg interface{}, path string, module string) error

type IApplicationTransportREST interface {
	Run(ctx context.Context) error
}
type IApplicationTransportGraphQL interface {
	Run(ctx context.Context) error
}
type IApplicationTransportWatermill interface {
	Run(ctx context.Context) error
}

type ApplicationContext struct {
	ConfigPath  string
	ServiceName string
	AgentListen AgentListen
	ReadConfig  ReadConfig
}

type RESTRouter struct {
}

type GraphQLResolver struct {
	Product productresolver.Resolver
}

type GraphQLMiddleware struct {
}

type GraphQLDataloader struct {
	Product *productloader.Dataloader
}
