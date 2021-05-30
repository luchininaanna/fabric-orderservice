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
	Address   string
}

type OrderRepository interface {
	Add(o Order) error
	Delete(orderUuid uuid.UUID) error
}
