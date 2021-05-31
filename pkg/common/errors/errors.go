package errors

import "errors"

var InternalError = errors.New("InternalServerError")
var InvalidArgumentError = errors.New("InvalidArgumentError")
var OrderNotExistError = errors.New("OrderNotExistError")
var OrderAlreadyDeletedError = errors.New("OrderAlreadyDeletedError")
var InvalidItemQuantityError = errors.New("InvalidItemQuantityError")
var EmptyOrderError = errors.New("EmptyOrderError")
var InvalidOrderCostError = errors.New("InvalidOrderCostError")
var EmptyOrderAddressError = errors.New("EmptyOrderAddressError")
