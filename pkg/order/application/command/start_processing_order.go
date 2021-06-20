package command

import (
	"github.com/google/uuid"
	appErrors "orderservice/pkg/order/application/errors"
	"orderservice/pkg/order/model"
)

type StartProcessingOrderCommand struct {
	ID uuid.UUID
}

type startProcessingOrderCommandHandler struct {
	unitOfWork UnitOfWork
}

type StartProcessingOrderCommandHandler interface {
	Handle(command StartProcessingOrderCommand) error
}

func NewStartProcessingOrderCommandHandler(unitOfWork UnitOfWork) StartProcessingOrderCommandHandler {
	return &startProcessingOrderCommandHandler{unitOfWork}
}

func (h *startProcessingOrderCommandHandler) Handle(c StartProcessingOrderCommand) error {
	err := h.unitOfWork.Execute(func(rp model.OrderRepository) error {
		order, err := rp.Get(c.ID)
		if err != nil {
			return appErrors.OrderNotExistError
		}

		err = order.StartProcessing()
		if err != nil {
			return err
		}

		return rp.Store(*order)
	})

	return err
}
