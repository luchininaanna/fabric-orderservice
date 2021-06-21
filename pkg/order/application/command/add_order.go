package command

import (
	"github.com/google/uuid"
	"math/rand"
	"orderservice/pkg/order/model"
	"time"
)

type OrderItem struct {
	ID       uuid.UUID
	Quantity float32
}

type AddOrderCommand struct {
	Items   []OrderItem
	Address string
}

type addOrderCommandHandler struct {
	unitOfWork UnitOfWork
}

type AddOrderCommandHandler interface {
	Handle(command AddOrderCommand) (*uuid.UUID, error)
}

func NewAddOrderCommandHandler(unitOfWork UnitOfWork) AddOrderCommandHandler {
	return &addOrderCommandHandler{unitOfWork}
}

func (h *addOrderCommandHandler) Handle(c AddOrderCommand) (*uuid.UUID, error) {
	var orderId *uuid.UUID
	err := h.unitOfWork.Execute(func(rp model.OrderRepository) error {
		//TODO: проверить, что в заказе указаны существующие items (используя второй сервис)
		// надо делать вне транзакции и лока
		cost := rand.Float32() //TODO: получить стоимость заказа (используя второй сервис)

		orderItemDtoList := make([]model.OrderItemDto, 0)
		for _, item := range c.Items {
			orderItemDto := model.OrderItemDto{
				ID:       item.ID,
				Quantity: item.Quantity,
			}

			orderItemDtoList = append(orderItemDtoList, orderItemDto)
		}

		order, err := model.NewOrder(uuid.New(), orderItemDtoList, time.Now(), cost, model.OrderStatusOrderCreated, c.Address)
		if err != nil {
			return err
		}

		orderId = &order.ID
		return rp.Store(order)
	})

	return orderId, err
}
