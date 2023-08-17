package graph

import "github.com/k-vanio/simple-example-of-clean-architecture/internal/usecases"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CreateOrderUseCase usecases.CreateOrderUseCase
}
