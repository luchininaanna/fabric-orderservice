package command

import (
	"github.com/google/uuid"
	"orderservice/pkg/order/application/errors"
	"testing"
)

func TestDeleteNotExistentOrder(t *testing.T) {
	uow := &mockUnitOfWork{}
	h := deleteOrderCommandHandler{uow}
	err := h.Handle(DeleteOrderCommand{
		uuid.New().String(),
	})

	if err != errors.OrderNotExistError {
		t.Error("Delete not existent order")
	}
}
