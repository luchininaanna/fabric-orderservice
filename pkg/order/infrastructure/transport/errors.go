package transport

import (
	"errors"
	log "github.com/sirupsen/logrus"
	commonErrors "orderservice/pkg/common/errors"
	appErrors "orderservice/pkg/order/application/errors"
)

var EmptyItemListError = errors.New("order: empty item list")
var InvalidOrderStatusError = errors.New("order: invalid order status")
var OrderContainsNonExistentItemError = errors.New("order: order contains non-existent item")
var InternalError = commonErrors.InternalError

func WrapError(err error) error {
	switch err {
	case nil:
		return nil
	case appErrors.EmptyItemListError:
		return EmptyItemListError
	case appErrors.OrderContainsNonExistentItemError:
		return OrderContainsNonExistentItemError
	case commonErrors.InternalError:
		return InternalError
	default:
		log.Error(err)
		return InternalError
	}
}
