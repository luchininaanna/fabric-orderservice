package command

import (
	"orderservice/pkg/order/application/errors"
	"testing"
)

func TestAddOrderWithEmptyItemList(t *testing.T) {
	uow := &mockUnitOfWork{}
	h := addOrderCommandHandler{uow}
	_, err := h.Handle(AddOrderCommand{
		[]OrderItem{},
		"",
	})

	if err != errors.EmptyItemListError {
		t.Error("Create order with empty items list")
	}
}

func TestAddOrderWithNonExistentItem(t *testing.T) {
	//TODO:: добавить тест для OrderContainsNonExistentItemError
}
