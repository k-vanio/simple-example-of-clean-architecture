package graph

import (
	"context"

	"github.com/k-vanio/simple-example-of-clean-architecture/internal/infra/graph/model"
)

// ListOrders is the resolver for the ListOrders field.
func (r *queryResolver) ListOrders(ctx context.Context) ([]*model.Order, error) {
	rows, err := r.OrderUseCase.List()
	if err != nil {
		return nil, err
	}

	orders := []*model.Order{}

	for _, row := range rows {
		orders = append(orders, &model.Order{
			ID:         row.ID,
			Price:      row.Price,
			Tax:        row.Tax,
			FinalPrice: row.FinalPrice,
		})
	}

	return orders, nil
}
