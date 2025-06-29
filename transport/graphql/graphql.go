package transportgraphql

import (
	"context"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"

	"gobase/config"
	"gobase/di/registry"
	"gobase/graphql"
	graphqlgen "gobase/graphql/generated"
	middlewaregraphql "gobase/internal/pkg/middleware/graphql"
)

type TransportModule struct {
	graphQLResolver      registry.GraphQLResolver
	middlewareDataloader middlewaregraphql.Dataloader
	middlewareOtel       middlewaregraphql.Otel
	config               *config.MainConfig
}

type TransportOpts struct {
	GraphQLResolver      registry.GraphQLResolver
	MiddlewareDataloader middlewaregraphql.Dataloader
	MiddlewareOtel       middlewaregraphql.Otel
	Config               *config.MainConfig
}

const defaultPort = "8181"

func NewTransport(opts TransportOpts) (registry.IApplicationTransportGraphQL, registry.CleanupFunc) {
	transportModule := &TransportModule{
		graphQLResolver:      opts.GraphQLResolver,
		middlewareDataloader: opts.MiddlewareDataloader,
		middlewareOtel:       opts.MiddlewareOtel,
		config:               opts.Config,
	}

	return transportModule, transportModule.Cleanup
}

func (m *TransportModule) Cleanup() {

}

func (m *TransportModule) Run(ctx context.Context) error {
	port := strconv.Itoa(m.config.Server.GraphQL.Port)
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graphqlgen.NewExecutableSchema(graphqlgen.Config{Resolvers: &graphql.Resolver{
		GraphQLResolver: m.graphQLResolver,
	}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	m.middlewareOtel(srv)

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", m.middlewareDataloader(srv))

	err := http.ListenAndServe(":"+port, nil)

	return err
}
