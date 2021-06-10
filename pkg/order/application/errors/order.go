package errors

import "errors"

var OrderNotExistError = errors.New("order not exist")
var OrderContainsNonExistentItemError = errors.New("order contains non-existent item")
