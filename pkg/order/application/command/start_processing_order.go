package command

import (
	"github.com/google/uuid"
	"orderservice/pkg/common/errors"
	appErrors "orderservice/pkg/order/application/errors"
	"orderservice/pkg/order/model"
)

type StartProcessingOrderCommand struct {
	ID string
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
		orderUuid, err := uuid.Parse(c.ID)
		if err != nil {
			return errors.InvalidArgumentError
		}

		order, err := rp.Get(orderUuid)
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
