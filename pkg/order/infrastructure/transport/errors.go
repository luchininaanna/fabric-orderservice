package transport

import (
	"errors"
	log "github.com/sirupsen/logrus"
	commonErrors "orderservice/pkg/common/errors"
	appErrors "orderservice/pkg/order/application/errors"
	modelErrors "orderservice/pkg/order/model"
)

var InternalError = commonErrors.InternalError
var OrderContainsNonExistentItemError = errors.New("order: order contains non-existent item")
var OrderNotExistError = errors.New("order: order not exist")
var OrderAlreadyClosedError = errors.New("order: order already closed")
var OrderAlreadyInProcessError = errors.New("order: order already processing")
var InvalidItemQuantityError = errors.New("order: invalid item quantity")
var EmptyOrderError = errors.New("order: empty order")
var InvalidOrderCostError = errors.New("order: invalid order cost")
var EmptyOrderAddressError = errors.New("order: empty order address")
var InvalidOrderStatusError = errors.New("order: invalid order status")
var OrderWithEmptyItemListError = errors.New("order: order with empty item list")

func WrapError(err error) error {
	switch err {
	case nil:
		return nil
	case commonErrors.InternalError:
		return InternalError
	case appErrors.OrderNotExistError:
		return OrderNotExistError
	case appErrors.OrderContainsNonExistentItemError:
		return OrderContainsNonExistentItemError
	case modelErrors.OrderNotExistError:
		return OrderNotExistError
	case modelErrors.OrderAlreadyClosedError:
		return OrderAlreadyClosedError
	case modelErrors.OrderAlreadyInProcessError:
		return OrderAlreadyInProcessError
	case modelErrors.InvalidItemQuantityError:
		return InvalidItemQuantityError
	case modelErrors.EmptyOrderError:
		return EmptyOrderError
	case modelErrors.InvalidOrderCostError:
		return InvalidOrderCostError
	case modelErrors.EmptyOrderAddressError:
		return EmptyOrderAddressError
	case InvalidOrderStatusError:
		return InvalidOrderStatusError
	case OrderWithEmptyItemListError:
		return OrderWithEmptyItemListError
	default:
		log.Error(err)
		return InternalError
	}
}
