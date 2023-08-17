package usecases

import (
	"github.com/k-vanio/simple-example-of-clean-architecture/internal/entity"
	"github.com/k-vanio/simple-example-of-clean-architecture/pkg/events"
)

type OrderInputDTO struct {
	ID    string  `json:"id" example:"any name"`
	Price float64 `json:"price" example:"0.70"`
	Tax   float64 `json:"tax" example:"0.05"`
}

type OrderOutputDTO struct {
	ID         string  `json:"id" example:"any name"`
	Price      float64 `json:"price" example:"0.70"`
	Tax        float64 `json:"tax" example:"0.05"`
	FinalPrice float64 `json:"final_price" example:"0.75"`
}

type OrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *OrderUseCase {
	return &OrderUseCase{
		OrderRepository: OrderRepository,
		OrderCreated:    OrderCreated,
		EventDispatcher: EventDispatcher,
	}
}

func (c *OrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}
	order.CalculateFinalPrice()
	if err := c.OrderRepository.Save(&order); err != nil {
		return OrderOutputDTO{}, err
	}

	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.Price + order.Tax,
	}

	c.OrderCreated.SetPayload(dto)
	c.EventDispatcher.Dispatch(c.OrderCreated)

	return dto, nil
}

func (c *OrderUseCase) List() ([]*OrderOutputDTO, error) {
	orders := []*OrderOutputDTO{}

	rows, err := c.OrderRepository.List()
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		orders = append(orders, &OrderOutputDTO{
			ID:         row.ID,
			Price:      row.Price,
			Tax:        row.Tax,
			FinalPrice: row.FinalPrice,
		})
	}

	return orders, nil
}
