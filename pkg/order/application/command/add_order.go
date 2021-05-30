package command

import (
	"github.com/google/uuid"
	"orderservice/pkg/order/application/errors"
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

		if len(c.Items) == 0 {
			return errors.InvalidItemsListError
		}

		cost := 78 //TODO:получить стоимость заказа

		orderItems := make([]model.OrderItem, 0)
		for _, item := range c.Items {
			itemUuid, err := uuid.Parse(item.ID)
			if err != nil {
				return err
			}
			orderItems = append(orderItems, model.OrderItem{ID: itemUuid, Quantity: item.Quantity})
		}

		id := uuid.New()
		orderId = &id
		return rp.Add(model.Order{
			ID:        id,
			Items:     orderItems,
			CreatedAt: time.Now(),
			Cost:      cost,
			Address:   c.Address,
		})
	})

	return orderId, err
}
