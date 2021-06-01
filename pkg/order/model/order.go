package model

import (
	"github.com/google/uuid"
	"time"
)

type OrderItem struct {
	ID       uuid.UUID
	Quantity int
}

type Order struct {
	ID        uuid.UUID
	Items     []OrderItem
	CreatedAt time.Time
	Cost      int
	Status    int
	Address   string
}

type OrderRepository interface {
	Store(o Order) error
	Get(orderUuid uuid.UUID) (*Order, error)
}

func NewOrder(orderUuid uuid.UUID, items []OrderItem, createdAt time.Time, cost int, status int, address string) (Order, error) {

	if len(items) == 0 {
		return Order{}, EmptyOrderError
	}

	if cost <= 0 {
		return Order{}, InvalidOrderCostError
	}

	if address == "" {
		return Order{}, EmptyOrderAddressError
	}

	return Order{
		orderUuid,
		items,
		createdAt,
		cost,
		status,
		address,
	}, nil
}

func NewOrderItem(itemUuid uuid.UUID, quantity int) (OrderItem, error) {

	if quantity <= 0 {
		return OrderItem{}, InvalidItemQuantityError
	}

	return OrderItem{
		itemUuid,
		quantity,
	}, nil
}

func (o *Order) Delete() error {
	if o.Status != OrderStatusOrderCanceled {
		o.Status = OrderStatusOrderCanceled
		return nil
	}

	return OrderAlreadyDeletedError
}
