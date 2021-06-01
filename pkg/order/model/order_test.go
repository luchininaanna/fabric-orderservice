package model

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

var mockOrderItems = []OrderItem{{uuid.New(), 5}}

func TestCreateOrderWithEmptyItemList(t *testing.T) {
	_, err := NewOrder(uuid.New(), []OrderItem{}, time.Now(), 77, 0, "Address")
	if err != EmptyOrderError {
		t.Error("Create order with empty item list")
	}
}

func TestCreateOrderWithInvalidCost(t *testing.T) {
	_, err := NewOrder(uuid.New(), mockOrderItems, time.Now(), 0, 0, "Address")
	if err != InvalidOrderCostError {
		t.Error("Create order with invalid cost")
	}
}

func TestCreateOrderWithEmptyAddress(t *testing.T) {
	_, err := NewOrder(uuid.New(), mockOrderItems, time.Now(), 77, 0, "")
	if err != EmptyOrderAddressError {
		t.Error("Create order with empty address")
	}
}

func TestCreateCorrectOrder(t *testing.T) {
	_, err := NewOrder(uuid.New(), mockOrderItems, time.Now(), 77, 0, "Address")
	if err != nil {
		t.Error("Create correct order with error")
	}
}

func TestCreateOrderItemWithBelowZeroQuantity(t *testing.T) {
	_, err := NewOrderItem(uuid.New(), -3)
	if err != InvalidItemQuantityError {
		t.Error("Create order item with below zero quantity")
	}
}

func TestCreateOrderItemWithZeroQuantity(t *testing.T) {
	_, err := NewOrderItem(uuid.New(), 0)
	if err != InvalidItemQuantityError {
		t.Error("Create order item with zero quantity")
	}
}

func TestCreateCorrectOrderItem(t *testing.T) {
	_, err := NewOrderItem(uuid.New(), 5)
	if err != nil {
		t.Error("Create correct order item with error")
	}
}

func TestDeleteOrder(t *testing.T) {
	order, err := NewOrder(uuid.New(), mockOrderItems, time.Now(), 77, 0, "Address")
	if err != nil {
		t.Error("Create correct order with error")
		return
	}

	err = order.Delete()
	if err == nil {
		t.Error("Delete order successfully")
	}
}

func TestDeleteAlreadyDeletedOrder(t *testing.T) {
	order, err := NewOrder(uuid.New(), mockOrderItems, time.Now(), 77, 3, "Address")
	if err != nil {
		t.Error("Create correct order with error")
		return
	}

	err = order.Delete()
	if err != OrderAlreadyDeletedError {
		t.Error("Delete already deleted order")
	}
}
