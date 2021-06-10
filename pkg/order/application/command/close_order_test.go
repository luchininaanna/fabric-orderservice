package command

import (
	"github.com/google/uuid"
	"orderservice/pkg/order/application/errors"
	"testing"
)

func TestCloseNotExistentOrder(t *testing.T) {
	uow := &mockUnitOfWork{}
	h := closeOrderCommandHandler{uow}
	err := h.Handle(CloseOrderCommand{
		uuid.New().String(),
	})

	if err != errors.OrderNotExistError {
		t.Error("Close not existent order")
	}
}
