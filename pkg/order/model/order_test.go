package model

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

var mockOrderItems = []OrderItemDto{{uuid.New(), 5}}

func TestCreateOrderWithEmptyItemList(t *testing.T) {
	_, err := NewOrder(uuid.New(), []OrderItemDto{}, time.Now(), 77, 0, "Address")
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
	_, err := newOrderItem(uuid.New(), -3)
	if err != InvalidItemQuantityError {
		t.Error("Create order item with below zero quantity")
	}
}

func TestCreateOrderItemWithZeroQuantity(t *testing.T) {
	_, err := newOrderItem(uuid.New(), 0)
	if err != InvalidItemQuantityError {
		t.Error("Create order item with zero quantity")
	}
}

func TestCreateCorrectOrderItem(t *testing.T) {
	_, err := newOrderItem(uuid.New(), 5)
	if err != nil {
		t.Error("Create correct order item with error")
	}
}

func TestCloseOrder(t *testing.T) {
	order, err := NewOrder(uuid.New(), mockOrderItems, time.Now(), 77, 0, "Address")
	if err != nil {
		t.Error("Create correct order with error")
		return
	}

	err = order.Close()
	if err != nil {
		t.Error("Close correct order with error")
	}
}

func TestCloseAlreadyClosedOrder(t *testing.T) {
	order, err := NewOrder(uuid.New(), mockOrderItems, time.Now(), 77, 3, "Address")
	if err != nil {
		t.Error("Create correct order with error")
		return
	}

	err = order.Close()
	if err != OrderAlreadyClosedError {
		t.Error("Close already closed order")
	}
}

func TestStartProcessingOrder(t *testing.T) {
	order, err := NewOrder(uuid.New(), mockOrderItems, time.Now(), 77, 0, "Address")
	if err != nil {
		t.Error("Create correct order with error")
		return
	}

	err = order.StartProcessing()
	if err != nil {
		t.Error("Start processing correct order with error")
	}
}

func TestStartProcessingAlreadyProcessingOrder(t *testing.T) {
	order, err := NewOrder(uuid.New(), mockOrderItems, time.Now(), 77, 3, "Address")
	if err != nil {
		t.Error("Create correct order with error")
		return
	}

	err = order.StartProcessing()
	if err != OrderAlreadyInProcessError {
		t.Error("Start processing already processing order")
	}
}
