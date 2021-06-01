package command

import (
	"github.com/google/uuid"
	"orderservice/pkg/common/errors"
	appErrors "orderservice/pkg/order/application/errors"
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
		//TODO:: сделать так, чтобы err возвращалась не nil, когда order не найден
		if err != nil {
			return appErrors.OrderNotExistError
		}
		if order == nil {
			return appErrors.OrderNotExistError
		}

		return order.Delete()
	})

	return err
}
