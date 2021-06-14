package model

import (
	"github.com/google/uuid"
	"time"
)

type OrderItemDto struct {
	ID       uuid.UUID
	Quantity int
}

type OrderItem struct {
	ID       uuid.UUID
	Quantity int
}

type Order struct {
	ID         uuid.UUID
	Items      []OrderItem
	CreatedAt  time.Time
	CanceledAt *time.Time
	Cost       int
	Status     int
	Address    string
}

type OrderRepository interface {
	Store(o Order) error
	Get(orderUuid uuid.UUID) (*Order, error)
}

func NewOrder(orderUuid uuid.UUID, items []OrderItemDto, createdAt time.Time, cost int, status int, address string) (Order, error) {

	if len(items) == 0 {
		return Order{}, EmptyOrderError
	}

	if cost <= 0 {
		return Order{}, InvalidOrderCostError
	}

	if address == "" {
		return Order{}, EmptyOrderAddressError
	}

	orderItems, err := getOrderItems(items)
	if err != nil {
		return Order{}, err
	}

	return Order{
		orderUuid,
		orderItems,
		createdAt,
		nil,
		cost,
		status,
		address,
	}, nil
}

func getOrderItems(items []OrderItemDto) ([]OrderItem, error) {
	orderItems := make([]OrderItem, 0)
	for _, item := range items {
		orderItem, err := newOrderItem(item.ID, item.Quantity)
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, orderItem)
	}

	return orderItems, nil
}

func newOrderItem(itemUuid uuid.UUID, quantity int) (OrderItem, error) {

	if quantity <= 0 {
		return OrderItem{}, InvalidItemQuantityError
	}

	return OrderItem{
		itemUuid,
		quantity,
	}, nil
}

func (o *Order) Close() error {
	if o.Status != OrderStatusOrderClosed {
		o.Status = OrderStatusOrderClosed
		now := time.Now()
		o.CanceledAt = &now
		return nil
	}
	return OrderAlreadyClosedError
}

func (o *Order) StartProcessing() error {
	if o.Status != OrderStatusOrderInProcess {
		o.Status = OrderStatusOrderInProcess
		return nil
	}
	return OrderAlreadyInProcessError
}
