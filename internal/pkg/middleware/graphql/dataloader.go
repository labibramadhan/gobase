package middlewaregraphql

import (
	"context"
	"net/http"

	"gobase/di/registry"
)

type contextKey string

const loadersKey = contextKey("dataloaders")

type Dataloader = func(next http.Handler) http.Handler

func NewDataloader(loader registry.GraphQLDataloader) Dataloader {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loadersKey, loader)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// For returns the dataloaders for the current request.
func For(ctx context.Context) registry.GraphQLDataloader {
	return ctx.Value(loadersKey).(registry.GraphQLDataloader)
}
