package command

import (
	"github.com/google/uuid"
	appErrors "orderservice/pkg/order/application/errors"
	"orderservice/pkg/order/model"
)

type CloseOrderCommand struct {
	ID uuid.UUID
}

type closeOrderCommandHandler struct {
	unitOfWork UnitOfWork
}

type CloseOrderCommandHandler interface {
	Handle(command CloseOrderCommand) error
}

func NewCloseOrderCommandHandler(unitOfWork UnitOfWork) CloseOrderCommandHandler {
	return &closeOrderCommandHandler{unitOfWork}
}

func (h *closeOrderCommandHandler) Handle(c CloseOrderCommand) error {
	err := h.unitOfWork.Execute(func(rp model.OrderRepository) error {
		order, err := rp.Get(c.ID)
		if err != nil {
			return appErrors.OrderNotExistError
		}

		err = order.Close()
		if err != nil {
			return err
		}

		return rp.Store(*order)
	})

	return err
}
