package command

import (
	"github.com/google/uuid"
	"orderservice/pkg/order/model"
)

type mockUnitOfWork struct {
	orders map[string]model.Order
}

func (m *mockUnitOfWork) Execute(f func(rp model.OrderRepository) error) error {
	return f(m)
}

func (m *mockUnitOfWork) Store(o model.Order) error {
	if m.orders == nil {
		m.orders = make(map[string]model.Order)
	}

	m.orders[o.ID.String()] = o
	return nil
}

func (m *mockUnitOfWork) Get(orderUuid uuid.UUID) (*model.Order, error) {
	if m.orders == nil {
		m.orders = make(map[string]model.Order)
	}

	o, ok := m.orders[orderUuid.String()]
	if ok {
		return &o, nil
	}

	return nil, nil
}

func (m *mockUnitOfWork) Delete(orderUuid uuid.UUID) error {
	if m.orders == nil {
		m.orders = make(map[string]model.Order)
	}

	_, ok := m.orders[orderUuid.String()]
	if ok {
		delete(m.orders, orderUuid.String())
	}

	return nil
}
