package command

import (
	"github.com/google/uuid"
	"orderservice/pkg/common/errors"
	"orderservice/pkg/order/model"
)

type DeleteOrderCommand struct {
	ID string
}

type deleteOrderCommandHandler struct {
	unitOfWork UnitOfWork
}

type DeleteOrderCommandHandler interface {
	Handle(command DeleteOrderCommand) error
}

func NewDeleteOrderCommandHandler(unitOfWork UnitOfWork) DeleteOrderCommandHandler {
	return &deleteOrderCommandHandler{unitOfWork}
}

func (h *deleteOrderCommandHandler) Handle(c DeleteOrderCommand) error {
	err := h.unitOfWork.Execute(func(rp model.OrderRepository) error {
		orderUuid, err := uuid.Parse(c.ID)
		if err != nil {
			return errors.InvalidArgumentError
		}
		order, err := rp.Get(orderUuid)
		if err != nil {
			return errors.OrderNotExistError
		}
		return order.Delete()
	})

	return err
}
