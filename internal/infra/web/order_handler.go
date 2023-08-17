package web

import (
	"encoding/json"
	"net/http"

	_ "github.com/k-vanio/simple-example-of-clean-architecture/docs"
	"github.com/k-vanio/simple-example-of-clean-architecture/internal/entity"
	"github.com/k-vanio/simple-example-of-clean-architecture/internal/usecases"
	"github.com/k-vanio/simple-example-of-clean-architecture/pkg/events"
	"github.com/k-vanio/simple-example-of-clean-architecture/pkg/response"
)

type Error struct {
	Err string `json:"error"`
}

type InputOrder usecases.OrderInputDTO
type OutOrder usecases.OrderOutputDTO

type WebOrderHandler struct {
	EventDispatcher   events.EventDispatcherInterface
	OrderRepository   entity.OrderRepositoryInterface
	OrderCreatedEvent events.EventInterface
}

func NewWebOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreatedEvent events.EventInterface,
) *WebOrderHandler {
	return &WebOrderHandler{
		EventDispatcher:   EventDispatcher,
		OrderRepository:   OrderRepository,
		OrderCreatedEvent: OrderCreatedEvent,
	}
}

// @Param request body InputOrder true "query params"
// @Accept json
// @Produce json
// @failure 400 {object} Error
// @failure 500 {object} Error
// @Success 201 {object} OutOrder
// @Router /orders [post]
func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecases.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		response.Json(w, http.StatusBadRequest, &Error{Err: err.Error()})
		return
	}

	useCase := usecases.NewOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
	output, err := useCase.Execute(dto)
	if err != nil {
		response.Json(w, http.StatusInternalServerError, &Error{Err: err.Error()})
		return
	}

	response.Json(w, http.StatusCreated, output)
}

// @Accept json
// @Produce json
// @failure 500 {object} Error
// @Success 200 {object} []OutOrder
// @Router /orders [get]
func (h *WebOrderHandler) List(w http.ResponseWriter, r *http.Request) {

	useCase := usecases.NewOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
	output, err := useCase.List()
	if err != nil {
		response.Json(w, http.StatusInternalServerError, &Error{Err: err.Error()})
		return
	}

	if output == nil {
		response.Json(w, http.StatusOK, []int{})
		return
	}

	response.Json(w, http.StatusOK, output)
}
