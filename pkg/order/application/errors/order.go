package errors

import "errors"

var EmptyItemListError = errors.New("empty item list")
var OrderNotExistError = errors.New("order not exist")
var OrderContainsNonExistentItemError = errors.New("order contains non-existent item")
