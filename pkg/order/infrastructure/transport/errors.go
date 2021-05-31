package transport

import (
	"errors"
	log "github.com/sirupsen/logrus"
	commonErrors "orderservice/pkg/common/errors"
	appErrors "orderservice/pkg/order/application/errors"
)

var InvalidItemsListError = errors.New("order: empty item list")
var InvalidOrderStatusError = errors.New("order: invalid order status")
var InternalError = commonErrors.InternalError

func WrapError(err error) error {
	switch err {
	case nil:
		return nil
	case appErrors.InvalidItemsListError:
		return InvalidItemsListError
	case commonErrors.InternalError:
		return InternalError
	default:
		log.Error(err)
		return InternalError
	}
}
