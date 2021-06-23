package command

import (
	"github.com/google/uuid"
	"orderservice/pkg/order/application/adapter"
	"orderservice/pkg/order/application/errors"
	"orderservice/pkg/order/model"
	"time"
)

type OrderItem struct {
	ID       uuid.UUID
	Quantity float32
}

type AddOrderCommand struct {
	Items   []OrderItem
	Address string
}

type addOrderCommandHandler struct {
	unitOfWork UnitOfWork
	store      adapter.StoreAdapter
}

type AddOrderCommandHandler interface {
	Handle(command AddOrderCommand) (*uuid.UUID, error)
}

func NewAddOrderCommandHandler(unitOfWork UnitOfWork, store adapter.StoreAdapter) AddOrderCommandHandler {
	return &addOrderCommandHandler{unitOfWork, store}
}

func (h *addOrderCommandHandler) Handle(c AddOrderCommand) (*uuid.UUID, error) {
	var orderId *uuid.UUID
	cost, err := h.getTotalOrderCost(c)
	if err != nil {
		return nil, err
	}

	err = h.unitOfWork.Execute(func(rp model.OrderRepository) error {
		orderItemDtoList := make([]model.OrderItemDto, 0)
		for _, item := range c.Items {
			orderItemDto := model.OrderItemDto{
				ID:       item.ID,
				Quantity: item.Quantity,
			}

			orderItemDtoList = append(orderItemDtoList, orderItemDto)
		}

		order, err := model.NewOrder(uuid.New(), orderItemDtoList, time.Now(), *cost, model.OrderStatusOrderCreated, c.Address)
		if err != nil {
			return err
		}

		orderId = &order.ID
		return rp.Store(order)
	})
	if err != nil {
		return nil, err
	}

	return orderId, nil
}

func (h *addOrderCommandHandler) getTotalOrderCost(command AddOrderCommand) (*float32, error) {
	fabrics, err := h.store.GetFabrics()
	if err != nil {
		return nil, err
	}

	idToCostMap := h.getFabricIdToFabricCostMap(fabrics)
	var totalCost float32

	for _, fabricInOrder := range command.Items {
		fabricInOrderId := fabricInOrder.ID.String()
		cost, ok := idToCostMap[fabricInOrderId]
		if !ok {
			return nil, errors.OrderContainsNonExistentItemError
		}

		totalCost = totalCost + cost
	}
	return &totalCost, nil
}

func (h *addOrderCommandHandler) getFabricIdToFabricCostMap(fabrics []adapter.Fabric) map[string]float32 {
	idToCostMap := make(map[string]float32)

	for _, fabric := range fabrics {
		idToCostMap[fabric.Id] = fabric.Cost
	}

	return idToCostMap
}
