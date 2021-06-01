package model

import "errors"

var OrderNotExistError = errors.New("OrderNotExistError")
var OrderAlreadyDeletedError = errors.New("OrderAlreadyDeletedError")
var InvalidItemQuantityError = errors.New("InvalidItemQuantityError")
var EmptyOrderError = errors.New("EmptyOrderError")
var InvalidOrderCostError = errors.New("InvalidOrderCostError")
var EmptyOrderAddressError = errors.New("EmptyOrderAddressError")
