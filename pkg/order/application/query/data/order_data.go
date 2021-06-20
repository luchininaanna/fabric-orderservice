package data

import (
	"time"
)

type OrderData struct {
	ID         string
	OrderItems []OrderItemData
	CreatedAt  time.Time
	Cost       float32
	Status     int
	Address    string
}
