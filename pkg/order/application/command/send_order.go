package command

import (
	"github.com/google/uuid"
	"orderservice/pkg/order/model"
)

type SendOrderCommand struct {
	ID uuid.UUID
}

type sendOrderCommandHandler struct {
	unitOfWork UnitOfWork
}

type SendOrderCommandHandler interface {
	Handle(command SendOrderCommand) error
}

func NewSendOrderCommandHandler(unitOfWork UnitOfWork) SendOrderCommandHandler {
	return &sendOrderCommandHandler{unitOfWork}
}

func (h *sendOrderCommandHandler) Handle(c SendOrderCommand) error {
	err := h.unitOfWork.Execute(func(rp model.OrderRepository) error {
		order, err := rp.Get(c.ID)
		if err != nil {
			return err
		}

		err = order.Send()
		if err != nil {
			return err
		}

		return rp.Store(*order)
	})

	return err
}
