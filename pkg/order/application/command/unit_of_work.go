package command

import "orderservice/pkg/order/model"

type UnitOfWork interface {
	Execute(func(rp model.OrderRepository) error) error
}
