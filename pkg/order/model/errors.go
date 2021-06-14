package model

import "errors"

var OrderNotExistError = errors.New("OrderNotExistError")
var OrderAlreadyClosedError = errors.New("OrderAlreadyClosedError")
var OrderAlreadyInProcessError = errors.New("OrderAlreadyInProcessError")
var InvalidItemQuantityError = errors.New("InvalidItemQuantityError")
var EmptyOrderError = errors.New("EmptyOrderError")
var InvalidOrderCostError = errors.New("InvalidOrderCostError")
var EmptyOrderAddressError = errors.New("EmptyOrderAddressError")
