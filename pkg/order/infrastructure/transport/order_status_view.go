package transport

import (
	"orderservice/pkg/order/model"
)

const OrderCreated = "Created"
const OrderInProcess = "InProcess"
const OrderCanceled = "Canceled"
const OrderShipped = "Shipped"

func WrapOrderStatus(orderStatus int) (string, error) {
	switch orderStatus {
	case model.OrderStatusOrderCreated:
		return OrderCreated, nil
	case model.OrderStatusOrderInProcess:
		return OrderInProcess, nil
	case model.OrderStatusOrderClosed:
		return OrderCanceled, nil
	case model.OrderStatusOrderShipped:
		return OrderShipped, nil
	default:
		return "", InvalidOrderStatusError
	}
}
