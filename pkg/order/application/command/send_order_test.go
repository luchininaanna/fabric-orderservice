package command

import (
	"github.com/google/uuid"
	"orderservice/pkg/order/model"
	"testing"
)

func TestSendNotExistentOrder(t *testing.T) {
	uow := &mockUnitOfWork{}
	h := sendOrderCommandHandler{uow}
	err := h.Handle(SendOrderCommand{
		uuid.New(),
	})

	if err != model.OrderNotExistError {
		t.Error("Send not existent order")
	}
}
