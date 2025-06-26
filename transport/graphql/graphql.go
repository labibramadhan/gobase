package transportgraphql

import (
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"

	"gobase/di/registry"
	"gobase/graphql"
	graphqlgen "gobase/graphql/generated"
	middlewaregraphql "gobase/internal/pkg/middleware/graphql"
)

type TransportModule struct {
	graphQLResolver      registry.GraphQLResolver
	middlewareDataloader middlewaregraphql.Dataloader
}

const defaultPort = "8181"

func NewTransport(graphQLResolver registry.GraphQLResolver, middlewareDataloader middlewaregraphql.Dataloader) registry.IApplicationTransportGraphQL {
	return &TransportModule{
		graphQLResolver:      graphQLResolver,
		middlewareDataloader: middlewareDataloader,
	}
}

func (m *TransportModule) Run() (registry.CleanupFunc, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graphqlgen.NewExecutableSchema(graphqlgen.Config{Resolvers: &graphql.Resolver{
		GraphQLResolver: m.graphQLResolver,
	}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", m.middlewareDataloader(srv))

	err := http.ListenAndServe(":"+port, nil)

	return nil, err
}
