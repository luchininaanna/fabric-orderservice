package command

import (
	"github.com/google/uuid"
	"orderservice/pkg/order/model"
	"time"
)

type OrderItem struct {
	ID       string
	Quantity int
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
		cost := 78 //TODO: получить стоимость заказа (используя второй сервис)

		orderItemDtoList := make([]model.OrderItemDto, 0)
		for _, item := range c.Items {
			itemUuid, err := uuid.Parse(item.ID)
			if err != nil {
				return err
			}

			orderItemDto := model.OrderItemDto{
				ID:       itemUuid,
				Quantity: item.Quantity,
			}

			orderItemDtoList = append(orderItemDtoList, orderItemDto)
		}

		order, err := model.NewOrder(uuid.New(), orderItemDtoList, time.Now(), cost, model.OrderStatusOrderCreated, c.Address)
		if err != nil {
			return err
		}

		return rp.Store(order)
	})

	return orderId, err
}
