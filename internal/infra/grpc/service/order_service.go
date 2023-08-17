package service

import (
	"context"

	"github.com/k-vanio/simple-example-of-clean-architecture/internal/infra/grpc/pb"
	"github.com/k-vanio/simple-example-of-clean-architecture/internal/usecases"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	OrderUseCase usecases.OrderUseCase
}

func NewOrderService(createOrderUseCase usecases.OrderUseCase) *OrderService {
	return &OrderService{
		OrderUseCase: createOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	dto := usecases.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.OrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.OrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(context.Context, *pb.Empty) (*pb.ListOrderResponse, error) {
	rows, err := s.OrderUseCase.List()
	if err != nil {
		return nil, err
	}

	orders := &pb.ListOrderResponse{}

	for _, row := range rows {
		orders.Orders = append(orders.Orders, &pb.OrderResponse{
			Id:         row.ID,
			Price:      float32(row.Price),
			Tax:        float32(row.Tax),
			FinalPrice: float32(row.FinalPrice),
		})
	}

	return orders, nil
}
