package data

import (
	"time"
)

type OrderData struct {
	ID         string
	OrderItems []OrderItemData
	CreatedAt  time.Time
	Cost       int
	Status     int
	Address    string
}
