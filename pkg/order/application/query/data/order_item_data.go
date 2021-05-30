package data

import (
	"github.com/google/uuid"
)

type OrderItemData struct {
	ID       uuid.UUID
	Quantity int
}
