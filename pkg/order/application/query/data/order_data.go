package data

import (
	"github.com/google/uuid"
	"time"
)

type OrderData struct {
	ID         uuid.UUID
	OrderItems []OrderItemData
	CreatedAt  time.Time
	Cost       int
	Address    string
}
