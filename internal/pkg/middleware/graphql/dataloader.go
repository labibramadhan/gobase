package middlewaregraphql

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"

	"gobase/di/registry"
)

type contextKey string

const loadersKey = contextKey("dataloaders")

type Dataloader = func(srv *handler.Server) http.Handler

func NewDataloader(loader registry.GraphQLDataloader) Dataloader {
	return func(srv *handler.Server) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loadersKey, loader)
			srv.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// For returns the dataloaders for the current request.
func For(ctx context.Context) registry.GraphQLDataloader {
	return ctx.Value(loadersKey).(registry.GraphQLDataloader)
}
