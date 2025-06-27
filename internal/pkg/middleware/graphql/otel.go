package middlewaregraphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/ravilushqa/otelgqlgen"
)

type Otel = func(srv *handler.Server)

func NewOtel() Otel {
	return func(srv *handler.Server) {
		srv.Use(otelgqlgen.Middleware())
	}
}
