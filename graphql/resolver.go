package graphql

//go:generate go run github.com/99designs/gqlgen generate

import "gobase/di/registry"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	GraphQLResolver registry.GraphQLResolver
	Dataloader      registry.GraphQLDataloader
}
