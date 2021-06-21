package command

import (
	"github.com/google/uuid"
	"orderservice/pkg/order/model"
	"testing"
)

func TestCloseNotExistentOrder(t *testing.T) {
	uow := &mockUnitOfWork{}
	h := closeOrderCommandHandler{uow}
	err := h.Handle(CloseOrderCommand{
		uuid.New(),
	})

	if err != model.OrderNotExistError {
		t.Error("Close not existent order")
	}
}
