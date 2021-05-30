package query

import "orderservice/pkg/order/application/query/data"

type OrderQueryService interface {
	GetOrder(id string) (*data.OrderData, error)
}
