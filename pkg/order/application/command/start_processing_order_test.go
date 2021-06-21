package command

import (
	"github.com/google/uuid"
	"orderservice/pkg/order/model"
	"testing"
)

func TestStartProcessingNotExistentOrder(t *testing.T) {
	uow := &mockUnitOfWork{}
	h := startProcessingOrderCommandHandler{uow}
	err := h.Handle(StartProcessingOrderCommand{
		uuid.New(),
	})

	if err != model.OrderNotExistError {
		t.Error("Start processing not existent order")
	}
}
